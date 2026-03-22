package helper

import (
	"reflect"
	"testing"
)

func AssertEqual(t *testing.T, got, want interface{}) {
	t.Helper()
	if got != want {
		t.Errorf("expected %v, but got %v", want, got)
	}
}

func AssertNotEqual(t *testing.T, got, want interface{}) {
	t.Helper()
	if got == want {
		t.Errorf("expected %v, but got %v", want, got)
	}
}

func AssertNil(t *testing.T, got interface{}) {
	t.Helper()
	if got == nil {
		return
	}

	v := reflect.ValueOf(got)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		if v.IsNil() {
			return
		}
	}
	t.Errorf("expected nil, but got %v", got)
}

func AssertNotNil(t *testing.T, got interface{}) {
	t.Helper()
	if got != nil {
		return
	}
	t.Errorf("expected not nil, but got %v", got)
}
