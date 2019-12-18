package config

import (
	"fmt"
	"strconv"
	"sync"
)

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
	DefValue string

	// Path of the value, this is a dot separated path internally (i.e. Debug.Enabled)
	Path string

	// Value of the setting
	Value Value

	notifiers sync.Map
}

// Set the Value from the provided string
func (s *Setting) Set(v string) error {
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
	default:
		// TODO: see if we have a MarshalSetting implementation
		return fmt.Errorf("type %T not supported", s.Value)
	}

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
	case uint:
		return strconv.FormatUint(uint64(val), 10)
	case *uint:
		return strconv.FormatUint(uint64(*val), 10)
	default:
		// TODO: see if we have an UnmarshalSetting implementation
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
	default:
		// TODO: see if we have an Equals implementation
		return fmt.Sprintf("%v", val) == v
	}
}
