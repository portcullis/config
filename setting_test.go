package config

import (
	"fmt"
	"reflect"
	"testing"
)

type setTest struct {
	To          interface{}
	Initializer interface{}
}

func (s setTest) Value() interface{} {
	return reflect.Indirect(reflect.ValueOf(s.Initializer)).Interface()
}

func (s setTest) Equals(v interface{}) bool {
	return reflect.Indirect(reflect.ValueOf(s.Initializer)).Interface() == v
}

func newSetTest(to interface{}, v interface{}) setTest {
	p := reflect.New(reflect.TypeOf(v))
	p.Elem().Set(reflect.ValueOf(v))

	return setTest{
		To:          to,
		Initializer: p.Interface(),
	}
}

func TestSetting_Set(t *testing.T) {
	tests := []setTest{
		newSetTest("changed", "initial"),
		newSetTest(true, false),

		newSetTest(int(23), int(5)),
		newSetTest(int8(23), int8(5)),
		newSetTest(int16(23), int16(5)),
		newSetTest(int32(23), int32(5)),
		newSetTest(int64(23), int64(5)),

		newSetTest(uint(23), uint(5)),
		newSetTest(uint8(23), uint8(5)),
		newSetTest(uint16(23), uint16(5)),
		newSetTest(uint32(23), uint32(5)),
		newSetTest(uint64(23), uint64(5)),
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%T", test.To), func(t *testing.T) {
			s := &Setting{Value: test.Initializer}

			if err := s.Set(fmt.Sprintf("%v", test.To)); err != nil {
				t.Fatalf("Failed to set string value: %v", err)
			}

			if !test.Equals(test.To) {
				t.Errorf("Failed to update value: expected %v; got %v", test.To, test.Value())
			}
		})
	}
}
