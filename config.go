package config

import "io"

// Default configuration Set
var Default = &Set{}

// New will create a new setting with the specified name, value, and description in the Default Set. Name can not be empty, value can not be nil
func New(name string, value Value, description string) *Setting {
	return Default.Setting(name, value, description)
}

// Get a setting by name
func Get(name string) *Setting {
	return Default.Get(name)
}

// Update an existing setting by name. This is useful to populate from command line and/or environment, etc...
func Update(name, value string) (bool, error) {
	return Default.Update(name, value)
}

// Subset will return a child Set of this Set
func Subset(name string) *Set {
	return Default.Subset(name)
}

// Bind the Pointer to a Struct. This will take all of the fields and attempt to create settings from them. Any child structs will be set in a subset of the parent struct by name. All fields will be passed into the Set.Setting() function as pointers so that the Set.Set() function can write to the underlying value.
//
// Fields names can be overwritten with the `setting` field tag.
//
// Descriptions on settings can be set with teh `description` field tag.
//
// You can mask the Stringer of the setting (set it to output *****) by setting the field tag `mask:"true"`. This is really important to do to passwords/tokens/etc... to make sure they don't end up in logs.
//
// If a `flag` field tag exists, the `setting.Flag()` function will be called with the value and `flag.CommandLine``
func Bind(value interface{}) *Set {
	return Default.Bind(value)
}

// Notify when any of the settings in this set, or any child set is added or changed
func Notify(n Notifier) *NotifyHandle {
	return Default.Notify(n)
}

// Range over the settings in the entire Set
func Range(fn func(string, *Setting) bool) {
	Default.Range(fn)
}

// Dump the current settings to the specified io.Writer in a tab separated list
func Dump(w io.Writer) error {
	return Default.Dump(w)
}
