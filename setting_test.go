package config

import (
	"bytes"
	"flag"
	"fmt"
	"reflect"
	"testing"
	"time"
)

type setTest struct {
	To          interface{}
	Initializer interface{}
	CheckString string
}

func (s setTest) Value() interface{} {
	return reflect.Indirect(reflect.ValueOf(s.Initializer)).Interface()
}

func (s setTest) Equals(v interface{}) bool {
	return reflect.Indirect(reflect.ValueOf(s.Initializer)).Interface() == v
}

func newSetTest(to interface{}, v interface{}, s string) setTest {
	p := reflect.New(reflect.TypeOf(v))
	p.Elem().Set(reflect.ValueOf(v))

	return setTest{
		To:          to,
		Initializer: p.Interface(),
		CheckString: s,
	}
}

func TestSetting_Set(t *testing.T) {
	tests := []setTest{
		newSetTest("changed", "initial", "initial"),
		newSetTest(true, false, "False"),
		newSetTest(time.Minute*23, time.Minute, "1m"),

		newSetTest(int(23), int(5), "5"),
		newSetTest(int8(23), int8(5), "5"),
		newSetTest(int16(23), int16(5), "5"),
		newSetTest(int32(23), int32(5), "5"),
		newSetTest(int64(23), int64(5), "5"),

		newSetTest(uint(23), uint(5), "5"),
		newSetTest(uint8(23), uint8(5), "5"),
		newSetTest(uint16(23), uint16(5), "5"),
		newSetTest(uint32(23), uint32(5), "5"),
		newSetTest(uint64(23), uint64(5), "5"),

		newSetTest(float32(23), float32(5), "5"),
		newSetTest(float64(23), float64(5), "5"),

		// actually treated like a uint8, but we make sure it works
		newSetTest(byte(6), byte(5), "5"),
	}

	for _, test := range tests {
		testName := fmt.Sprintf("%T", test.To)

		// override the internal type resolutions to test type aliases
		if v, ok := test.To.(byte); ok && v == 6 {
			testName = "byte"
		}

		t.Run(testName, func(t *testing.T) {
			s := &Setting{Value: test.Initializer}

			// validates if the provided string matches the formatting of the string value
			if !s.Equals(test.CheckString) {
				t.Errorf("Failed to equality check supplied string value: expected %q; got %q", test.CheckString, s.String())
			}

			// make sure we can set the value from the provided string value
			if err := s.Set(test.CheckString); err != nil {
				t.Errorf("Failed to set from string value: %v", err)
			}

			// make sure we can set the value from the raw sprinted to string
			if err := s.Set(fmt.Sprintf("%v", test.To)); err != nil {
				t.Fatalf("Failed to set string value: %v", err)
			}

			// validate that the pointer was in fact changed to the new value
			if !test.Equals(test.To) {
				t.Errorf("Failed to update value: expected %v; got %v", test.To, test.Value())
			}

			// validate that we don't get a blank string back (could probably be a better test TBH)
			if s.String() == "" {
				t.Errorf("Failed to string value: got %q", s.String())
			}

			// validate that the fmt.sprintf matches the equality checker
			if !s.Equals(fmt.Sprintf("%v", test.To)) {
				t.Errorf("Failed to equality check string value: expected %q; got %q", fmt.Sprintf("%v", test.To), s.String())
			}
		})
	}
}

type customSetting struct {
	Value       []byte
	Marshaled   bool
	Unmarshaled bool
	Equaled     bool
}

func (cs *customSetting) UnmarshalSetting(v string) error {
	cs.Value = []byte(v)
	cs.Unmarshaled = true
	return nil
}

func (cs *customSetting) MarshalSetting() string {
	cs.Marshaled = true
	return string(cs.Value)
}

func (cs *customSetting) Equals(v string) bool {
	cs.Equaled = true
	return bytes.Equal(cs.Value, []byte(v))
}

func TestSetting_CustomType(t *testing.T) {
	cs := &customSetting{
		Value: []byte("hello"),
	}

	st := &Setting{Value: cs}

	if string(cs.Value) != st.String() {
		t.Errorf("Failed to get string value for custom type")
	}
	if !cs.Marshaled {
		t.Error("Custom object MarshalSetting not called")
	}

	newValue := "goodbye"

	if err := st.Set(newValue); err != nil {
		t.Fatalf("Failed to set string value for custom type: %v", err)
	}

	if !cs.Unmarshaled {
		t.Errorf("Custom object UnmarshalSetting not called")
	}

	if string(cs.Value) != newValue {
		t.Errorf("Failed to get updated string value for custom type: expected %q; got %q", newValue, string(cs.Value))
	}

	if !st.Equals(newValue) {
		t.Error("Failed to match equality on setting to set value")
	}

	if !cs.Equaled {
		t.Errorf("Custom object Equals not called")
	}
}

func TestSetting_Notify(t *testing.T) {
	name := "Test"
	value1 := "value1"
	value2 := "value2"

	st := &Setting{Name: name, Value: value1}

	notifyCalled := false
	nh := st.Notify(NotifyFunc(func(s *Setting) {
		if s.Name != name {
			t.Errorf("Notification Setting Name did not Match expected name; expected %q got %q", name, s.Name)
		}
		notifyCalled = true
	}))

	if err := st.Set(value1); err != nil {
		t.Fatalf("Failed to set value: %v", err)
	}

	if notifyCalled {
		t.Errorf("Notification unexpectingly called when value was the same")
	}
	notifyCalled = false

	if err := st.Set(value2); err != nil {
		t.Fatalf("Failed to set value: %v", err)
	}

	if !notifyCalled {
		t.Errorf("Notification did not execute as expected when value changed; current %q", st.String())
	}
	notifyCalled = false

	if err := nh.Close(); err != nil {
		t.Fatalf("Failed to close Notify Handle: %v", err)
	}

	if err := st.Set(value1); err != nil {
		t.Fatalf("Failed to set value: %v", err)
	}

	if notifyCalled {
		t.Errorf("Notification unexpectingly called after Notify Handler Closed")
	}

}

func TestSetting_FlagCompat(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	st := &Setting{Name: "debug", Description: "Sets debug mode", Value: false}
	st.Flag("debug", fs)

	if err := fs.Parse([]string{"-debug"}); err != nil {
		t.Errorf("Failed to set debug flag: %v", err)
	}

	if st.Value.(bool) != true {
		t.Errorf("Failed to set Setting from flag -debug")
	}

	if st.Type() != "bool" {
		t.Errorf("Failed to resolve type; expected %q got %q", "bool", st.Type())
	}
}
