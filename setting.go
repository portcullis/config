package config

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Marshaler is the interface implemented by types that can marshal themselves into a setting string.
type Marshaler interface {
	MarshalSetting() string
}

// Unmarshaler is the interface implemented by types that can unmarshal a string setting of themselves.
type Unmarshaler interface {
	UnmarshalSetting(string) error
}

// Equality is the interface implemented by type to validate equality of the supplied string to themselves.
type Equality interface {
	Equals(string) bool
}

// Value is an interface for interacting with the underlying configuration value
type Value interface{}

// Setting within the configuration containing a Value
type Setting struct {
	// Mask will overwrite the String function to return ***** to protect from logging
	Mask bool

	// Name of the value
	Name string

	// Description of this setting, useful for help text
	Description string

	// DefaultValue of the Setting as a string
	DefaultValue string

	// Path of the value, this is a dot separated path internally (i.e. Debug.Enabled)
	Path string

	// Value of the setting
	Value Value

	notifiers sync.Map
}

// IsDefault will return if the value matches the default value specified in Setting.DefaultValue
func (s *Setting) IsDefault() bool {
	return s.Equals(s.DefaultValue)
}

// Notify provides a callback interface to when a setting has changed via Setting.Set
func (s *Setting) Notify(n Notifier) *NotifyHandle {
	if n == nil {
		return &NotifyHandle{}
	}

	handle := &NotifyHandle{
		stopFunc: s.notifiers.Delete,
	}

	s.notifiers.Store(handle, n)

	return handle
}

// Set the Value from the provided string
func (s *Setting) Set(v string) error {
	same := s.Equals(v)

	switch val := s.Value.(type) {
	case string:
		s.Value = v
	case *string:
		*val = v
	case bool:
		pv, err := strconv.ParseBool(v)
		if err != nil {
			return fmt.Errorf("unable to cast value to boolean: %w", err)
		}
		s.Value = pv
	case *bool:
		pv, err := strconv.ParseBool(v)
		if err != nil {
			return fmt.Errorf("unable to cast value to boolean: %w", err)
		}
		*val = pv

	case int:
		pv, err := strconv.ParseInt(v, 0, strconv.IntSize)
		if err != nil {
			return fmt.Errorf("unable to cast value to int: %w", err)
		}
		s.Value = int(pv)
	case *int:
		pv, err := strconv.ParseInt(v, 0, strconv.IntSize)
		if err != nil {
			return fmt.Errorf("unable to cast value to int: %w", err)
		}
		*val = int(pv)
	case int8:
		pv, err := strconv.ParseInt(v, 0, 8)
		if err != nil {
			return fmt.Errorf("unable to cast value to int8: %w", err)
		}
		s.Value = int8(pv)
	case *int8:
		pv, err := strconv.ParseInt(v, 0, 8)
		if err != nil {
			return fmt.Errorf("unable to cast value to int8: %w", err)
		}
		*val = int8(pv)
	case int16:
		pv, err := strconv.ParseInt(v, 0, 16)
		if err != nil {
			return fmt.Errorf("unable to cast value to int16: %w", err)
		}
		s.Value = int16(pv)
	case *int16:
		pv, err := strconv.ParseInt(v, 0, 16)
		if err != nil {
			return fmt.Errorf("unable to cast value to int16: %w", err)
		}
		*val = int16(pv)
	case int32:
		pv, err := strconv.ParseInt(v, 0, 32)
		if err != nil {
			return fmt.Errorf("unable to cast value to int32: %w", err)
		}
		s.Value = int32(pv)
	case *int32:
		pv, err := strconv.ParseInt(v, 0, 32)
		if err != nil {
			return fmt.Errorf("unable to cast value to int32: %w", err)
		}
		*val = int32(pv)
	case int64:
		pv, err := strconv.ParseInt(v, 0, 64)
		if err != nil {
			return fmt.Errorf("unable to cast value to int64: %w", err)
		}
		s.Value = pv
	case *int64:
		pv, err := strconv.ParseInt(v, 0, 64)
		if err != nil {
			return fmt.Errorf("unable to cast value to int64: %w", err)
		}
		*val = pv

	case uint:
		pv, err := strconv.ParseUint(v, 0, strconv.IntSize)
		if err != nil {
			return fmt.Errorf("unable to case value to uint: %w", err)
		}
		s.Value = uint(pv)
	case *uint:
		pv, err := strconv.ParseUint(v, 0, strconv.IntSize)
		if err != nil {
			return fmt.Errorf("unable to case value to uint: %w", err)
		}
		*val = uint(pv)
	case uint8:
		pv, err := strconv.ParseUint(v, 0, 8)
		if err != nil {
			return fmt.Errorf("unable to cast value to uint8: %w", err)
		}
		s.Value = uint8(pv)
	case *uint8:
		pv, err := strconv.ParseUint(v, 0, 8)
		if err != nil {
			return fmt.Errorf("unable to cast value to uint8: %w", err)
		}
		*val = uint8(pv)
	case uint16:
		pv, err := strconv.ParseUint(v, 0, 16)
		if err != nil {
			return fmt.Errorf("unable to cast value to uint16: %w", err)
		}
		s.Value = uint16(pv)
	case *uint16:
		pv, err := strconv.ParseUint(v, 0, 16)
		if err != nil {
			return fmt.Errorf("unable to cast value to uint16: %w", err)
		}
		*val = uint16(pv)
	case uint32:
		pv, err := strconv.ParseUint(v, 0, 32)
		if err != nil {
			return fmt.Errorf("unable to cast value to uint32: %w", err)
		}
		s.Value = uint32(pv)
	case *uint32:
		pv, err := strconv.ParseUint(v, 0, 32)
		if err != nil {
			return fmt.Errorf("unable to cast value to uint32: %w", err)
		}
		*val = uint32(pv)
	case uint64:
		pv, err := strconv.ParseUint(v, 0, 64)
		if err != nil {
			return fmt.Errorf("unable to cast value to uint64: %w", err)
		}
		s.Value = pv
	case *uint64:
		pv, err := strconv.ParseUint(v, 0, 64)
		if err != nil {
			return fmt.Errorf("unable to cast value to uint64: %w", err)
		}
		*val = pv

	case float32:
		pv, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return fmt.Errorf("unable to cast value to float32: %w", err)
		}
		s.Value = float32(pv)
	case *float32:
		pv, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return fmt.Errorf("unable to cast value to float32: %w", err)
		}
		*val = float32(pv)
	case float64:
		pv, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return fmt.Errorf("unable to cast value to float64: %w", err)
		}
		s.Value = pv
	case *float64:
		pv, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return fmt.Errorf("unable to cast value to float64: %w", err)
		}
		*val = pv

	case time.Duration:
		pv, err := time.ParseDuration(v)
		if err != nil {
			return fmt.Errorf("unable to cast value to time.Duration: %w", err)
		}
		s.Value = pv
	case *time.Duration:
		pv, err := time.ParseDuration(v)
		if err != nil {
			return fmt.Errorf("unable to cast value to time.Duration: %w", err)
		}
		*val = pv

	default:
		if unmarshaler, ok := s.Value.(Unmarshaler); ok {
			if err := unmarshaler.UnmarshalSetting(v); err != nil {
				return fmt.Errorf("unable to marshal value to %T: %w", s.Value, err)
			}
		} else {
			return fmt.Errorf("type %T not supported", s.Value)
		}
	}

	// if same, then go ahead and exit the function and don't notify
	if same {
		return nil
	}

	// notify those of changed value
	s.notifiers.Range(func(key, val interface{}) bool {
		f, ok := val.(Notifier)
		if !ok || f == nil {
			s.notifiers.Delete(key)
			return true
		}

		f.Notify(s)

		return true
	})

	return nil
}

func (s *Setting) String() string {
	if s.Mask {
		return "*****"
	}

	switch val := s.Value.(type) {
	case string:
		return val
	case *string:
		return *val
	case bool:
		return strconv.FormatBool(val)
	case *bool:
		return strconv.FormatBool(*val)

	case int:
		return strconv.Itoa(val)
	case *int:
		return strconv.Itoa(*val)
	case int8:
		return strconv.FormatInt(int64(val), 10)
	case *int8:
		return strconv.FormatInt(int64(*val), 10)
	case int16:
		return strconv.FormatInt(int64(val), 10)
	case *int16:
		return strconv.FormatInt(int64(*val), 10)
	case int32:
		return strconv.FormatInt(int64(val), 10)
	case *int32:
		return strconv.FormatInt(int64(*val), 10)
	case int64:
		return strconv.FormatInt(val, 10)
	case *int64:
		return strconv.FormatInt(*val, 10)

	case uint:
		return strconv.FormatUint(uint64(val), 10)
	case *uint:
		return strconv.FormatUint(uint64(*val), 10)
	case uint8:
		return strconv.FormatUint(uint64(val), 10)
	case *uint8:
		return strconv.FormatUint(uint64(*val), 10)
	case uint16:
		return strconv.FormatUint(uint64(val), 10)
	case *uint16:
		return strconv.FormatUint(uint64(*val), 10)
	case uint32:
		return strconv.FormatUint(uint64(val), 10)
	case *uint32:
		return strconv.FormatUint(uint64(*val), 10)
	case uint64:
		return strconv.FormatUint(val, 10)
	case *uint64:
		return strconv.FormatUint(*val, 10)

	case float32:
		return strconv.FormatFloat(float64(val), 'g', -1, 32)
	case *float32:
		return strconv.FormatFloat(float64(*val), 'g', -1, 32)
	case float64:
		return strconv.FormatFloat(val, 'g', -1, 64)
	case *float64:
		return strconv.FormatFloat(*val, 'g', -1, 64)

	default:
		if marshaler, ok := s.Value.(Marshaler); ok {
			return marshaler.MarshalSetting()
		}

		return fmt.Sprintf("%v", val)
	}
}

// Equals will validate that the input string is the same as the current value using the internal parsing
func (s *Setting) Equals(v string) bool {
	switch val := s.Value.(type) {
	case string:
		return val == v
	case *string:
		return *val == v
	case bool:
		pv, err := strconv.ParseBool(v)
		if err != nil {
			return false
		}
		return val == pv
	case *bool:
		pv, err := strconv.ParseBool(v)
		if err != nil {
			return false
		}
		return *val == pv

	case int:
		pv, err := strconv.ParseInt(v, 0, strconv.IntSize)
		if err != nil {
			return false
		}
		return val == int(pv)
	case *int:
		pv, err := strconv.ParseInt(v, 0, strconv.IntSize)
		if err != nil {
			return false
		}
		return *val == int(pv)
	case int8:
		pv, err := strconv.ParseInt(v, 0, 8)
		if err != nil {
			return false
		}
		return val == int8(pv)
	case *int8:
		pv, err := strconv.ParseInt(v, 0, 8)
		if err != nil {
			return false
		}
		return *val == int8(pv)
	case int16:
		pv, err := strconv.ParseInt(v, 0, 16)
		if err != nil {
			return false
		}
		return val == int16(pv)
	case *int16:
		pv, err := strconv.ParseInt(v, 0, 16)
		if err != nil {
			return false
		}
		return *val == int16(pv)
	case int32:
		pv, err := strconv.ParseInt(v, 0, 32)
		if err != nil {
			return false
		}
		return val == int32(pv)
	case *int32:
		pv, err := strconv.ParseInt(v, 0, 32)
		if err != nil {
			return false
		}
		return *val == int32(pv)
	case int64:
		pv, err := strconv.ParseInt(v, 0, 64)
		if err != nil {
			return false
		}
		return val == pv
	case *int64:
		pv, err := strconv.ParseInt(v, 0, 64)
		if err != nil {
			return false
		}
		return *val == pv

	case uint:
		pv, err := strconv.ParseUint(v, 0, strconv.IntSize)
		if err != nil {
			return false
		}
		return val == uint(pv)
	case *uint:
		pv, err := strconv.ParseUint(v, 0, strconv.IntSize)
		if err != nil {
			return false
		}
		return *val == uint(pv)
	case uint8:
		pv, err := strconv.ParseUint(v, 0, 8)
		if err != nil {
			return false
		}
		return val == uint8(pv)
	case *uint8:
		pv, err := strconv.ParseUint(v, 0, 8)
		if err != nil {
			return false
		}
		return *val == uint8(pv)
	case uint16:
		pv, err := strconv.ParseUint(v, 0, 16)
		if err != nil {
			return false
		}
		return val == uint16(pv)
	case *uint16:
		pv, err := strconv.ParseUint(v, 0, 16)
		if err != nil {
			return false
		}
		return *val == uint16(pv)
	case uint32:
		pv, err := strconv.ParseUint(v, 0, 32)
		if err != nil {
			return false
		}
		return val == uint32(pv)
	case *uint32:
		pv, err := strconv.ParseUint(v, 0, 32)
		if err != nil {
			return false
		}
		return *val == uint32(pv)
	case uint64:
		pv, err := strconv.ParseUint(v, 0, 64)
		if err != nil {
			return false
		}
		return val == pv
	case *uint64:
		pv, err := strconv.ParseUint(v, 0, 64)
		if err != nil {
			return false
		}
		return *val == pv

	case float32:
		pv, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return false
		}
		return val == float32(pv)
	case *float32:
		pv, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return false
		}
		return *val == float32(pv)

	case float64:
		pv, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return false
		}
		return val == pv
	case *float64:
		pv, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return false
		}
		return *val == pv

	case time.Duration:
		pv, err := time.ParseDuration(v)
		if err != nil {
			return false
		}
		return val == pv
	case *time.Duration:
		pv, err := time.ParseDuration(v)
		if err != nil {
			return false
		}
		return *val == pv

	default:
		if equality, ok := s.Value.(Equality); ok {
			return equality.Equals(v)
		}

		return fmt.Sprintf("%v", val) == v
	}
}

// Type returns a string representation of the type, but omits the pointer prefix (*)
// This is provided to complete the interface for the github.com/spf13/pflag package
func (s *Setting) Type() string {
	return strings.TrimLeft(fmt.Sprintf("%T", s.Value), "*")
}

// IsBoolFlag is provided to help support boolean flags in the flag package (i.e. -debug rather than -debug=true)
func (s *Setting) IsBoolFlag() bool {
	switch s.Value.(type) {
	case *bool, bool:
		return true
	default:
		return false
	}
}

// Flag will register the current Setting as a command line flag in the supplied flag.FlagSet. When the supplied fs is nill, the flag.CommandLine is used
func (s *Setting) Flag(arg string, fs *flag.FlagSet) {
	if fs == nil {
		fs = flag.CommandLine
	}

	fs.Var(s, arg, s.Description)
}
