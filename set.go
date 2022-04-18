package config

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
	"text/tabwriter"
)

// Set defines a composite collection of configuration
type Set struct {
	name      string
	path      string
	root      *Set
	parent    *Set
	children  sync.Map
	settings  sync.Map
	notifiers sync.Map
}

// Get a setting by name
func (s *Set) Get(name string) *Setting {
	root := s.root
	if root == nil {
		root = s
	}

	if setting, found := root.settings.Load(strings.ToLower(name)); found {
		return setting.(*Setting)
	}

	path := fmt.Sprintf("%s.%s", s.path, name)
	if setting, found := root.settings.Load(strings.ToLower(path)); found {
		return setting.(*Setting)
	}

	return nil
}

// Set an existing setting by name. This is useful to populate from command line and/or environment, etc...
func (s *Set) Set(name, value string) (bool, error) {
	setting := s.Get(name)
	if setting == nil {
		return false, nil
	}

	return true, setting.Set(value)
}

// Subset will return a child Set of this Set
func (s *Set) Subset(name string) *Set {
	root := s.root
	if root == nil {
		root = s
	}

	subsetPath := fmt.Sprintf("%s.%s", s.path, name)
	if s.path == "" {
		subsetPath = name
	}

	if set, found := root.children.Load(strings.ToLower(subsetPath)); found {
		return set.(*Set)
	}

	set := &Set{
		name:   name,
		path:   subsetPath,
		root:   root,
		parent: s,
	}

	root.children.Store(strings.ToLower(subsetPath), set)

	return set
}

// Path of the Set, child Set's will have a dot separated path (root.child.child)
func (s *Set) Path() string {
	return s.path
}

// Name of the current set
func (s *Set) Name() string {
	return s.name
}

// Root set of the config
func (s *Set) Root() *Set {
	if s.root == nil {
		return s
	}

	return s.root
}

// Parent of the current set
func (s *Set) Parent() *Set {
	if s.parent == nil {
		return s
	}

	return s.parent
}

// Setting will create a new setting with the specified name, value, and description in the current Set. Name can not be empty, value can not be nil
func (s *Set) Setting(name string, value Value, description string) *Setting {
	if name == "" {
		panic("name can not be empty")
	}
	if value == nil {
		panic("value can not be nil")
	}

	root := s.root
	if root == nil {
		root = s
	}

	settingPath := fmt.Sprintf("%s.%s", s.path, name)
	if s.path == "" {
		settingPath = name
	}

	setting := &Setting{
		Name:        name,
		Description: description,
		Path:        settingPath,
		Value:       value,
	}

	// cheeky allows the underlying thing to actually map it properly
	setting.DefaultValue = setting.String()

	_, exists := root.settings.LoadOrStore(strings.ToLower(settingPath), setting)
	if exists {
		panic(fmt.Sprintf("setting %q already exists", settingPath))
	}

	// get notified when the setting changes - we won't stop notifications as long as it is a child, and since there is no remove.... we just discard the Close handler
	_ = setting.Notify(NotifyFunc(s.notifyChanged))

	// notify that we have added something (a change) after returning
	defer s.notifyChanged(setting)

	return setting
}

// Range over the settings in the entire Set
func (s *Set) Range(fn func(string, *Setting) bool) {
	root := s.root
	if root == nil {
		root = s
	}

	root.settings.Range(func(k, v any) bool {
		key := k.(string)
		setting := v.(*Setting)

		if !strings.HasPrefix(key, s.path) {
			return true
		}

		return fn(key, setting)
	})
}

// Bind the Pointer to a Struct. This will take all of the fields and attempt to create settings from them. Any child structs will be set in a subset of the parent struct by name. All fields will be passed into the Set.Setting() function as pointers so that the Set.Set() function can write to the underlying value.
//
// Fields names can be overwritten with the `setting` field tag.
//
// Descriptions on settings can be set with teh `description` field tag.
//
// You can mask the Stringer of the setting (set it to output *****) by setting the field tag `mask:"true"`. This is really important to do to passwords/tokens/etc... to make sure they don't end up in logs.
func (s *Set) Bind(value interface{}) {
	rvalue := reflect.ValueOf(value)

	if rvalue.Kind() != reflect.Ptr {
		panic("value must be a pointer value")
	}

	rvalue = rvalue.Elem()

	if rvalue.Kind() != reflect.Struct {
		panic("value must be a struct value")
	}

	for i := 0; i < rvalue.NumField(); i++ {
		fieldType := rvalue.Type().Field(i)
		fieldValue := rvalue.Field(i)

		if !fieldValue.CanSet() {
			continue
		}

		description := fieldType.Tag.Get("description")
		name := fieldType.Name
		masked := fieldType.Tag.Get("mask") == "true"
		if tagName := fieldType.Tag.Get("setting"); tagName != "" {
			name = tagName
		}

		if name == "-" {
			continue
		}

		switch rvalue.Field(i).Kind() {
		case reflect.Invalid, reflect.Chan, reflect.Func:
			// do nothing

		case reflect.Ptr:
			// if the thing is a pointer, then call this as a child
			s.Subset(name).Bind(fieldValue.Interface())

		case reflect.Struct:
			// if the thing is a struct, pass it through as a child
			s.Subset(name).Bind(fieldValue.Addr().Interface())

		default:
			// all other field types we pass in the pointer to the value as a setting so that it is "bound"
			setting := s.Setting(name, fieldValue.Addr().Interface(), description)
			setting.Mask = masked
		}
	}
}

// Dump the current settings to the specified io.Writer in a tab separated list
func (s *Set) Dump(w io.Writer) error {
	tw := tabwriter.NewWriter(w, 10, 10, 5, ' ', 0)

	fmt.Fprintln(tw, "Path\tType\tValue\tDescription")

	s.Range(func(path string, setting *Setting) bool {
		fmt.Fprintf(tw, "%s\t%T\t%q\t%s\t\n", setting.Path, setting.Value, setting.String(), setting.Description)
		return true
	})

	return tw.Flush()
}

// Notify when any of the settings in this set, or any child set is added or changed
func (s *Set) Notify(n Notifier) *NotifyHandle {
	if n == nil {
		return &NotifyHandle{}
	}

	handle := &NotifyHandle{
		stopFunc: s.notifiers.Delete,
	}

	s.notifiers.Store(handle, n)

	return handle
}

// notifyChanged is attached to all settings so that we can get notified of when they are added
func (s *Set) notifyChanged(setting *Setting) {
	s.notifiers.Range(func(k, v interface{}) bool {
		notifier := v.(Notifier)
		notifier.Notify(setting)
		return true
	})

	// call the parent to notify if they exist to propogate upward the notification
	if s.parent != nil {
		s.parent.notifyChanged(setting)
	}
}
