package di

import (
	"context"
	"reflect"
	"testing"
)

type TestType struct{}
type AnotherTestType struct{}

var (
	diContainer Container
	ctx 			  context.Context
)

func init() {
	ctx = context.Background()
	diContainer = NewDI(ctx)
}

func TestDI_IsEmpty(t *testing.T) {
	if !diContainer.IsEmpty() {
		t.Errorf("DI Container is not empty")
	}
}

func TestDI_IsEmpty_AfterMustSet(t *testing.T) {
	diContainer.Clear()

	diContainer.MustSet("test", &TestType{})

	if diContainer.IsEmpty() {
		t.Errorf("DI Container must not be empty after storing test object")
	}
}

func TestDI_IsEmpty_AfterClearAndMustSet(t *testing.T) {
	diContainer.Clear()
	diContainer.MustSet("test", &TestType{})

	if diContainer.IsEmpty() {
		t.Errorf("DI Container must not be empty after storing test object")
	}
}

func TestDI_IsEmpty_AfterMustSetAndClear(t *testing.T) {
	diContainer.Clear()

	diContainer.MustSet("test", &TestType{})
	diContainer.Clear()

	if !diContainer.IsEmpty() {
		t.Errorf("DI Container must be empty after clearing")
	}
}

func TestDI_Size(t *testing.T) {
	diContainer.Clear()
	diContainer.MustSet("test", &TestType{})

	if diContainer.Size() != 1 {
		t.Errorf("DI Container must have only one object")
	}
}

func TestDI_Size_AfterClear(t *testing.T) {
	diContainer.Clear()

	diContainer.MustSet("test", &TestType{})
	diContainer.Clear()

	if diContainer.Size() != 0 {
		t.Errorf("DI Container must have 0 objects after clearing")
	}
}

func TestDI_Set(t *testing.T) {
	diContainer.Clear()

	expected := &TestType{}
	if err := diContainer.Set("test", expected); err != nil {
		t.Error(err.Error())
	}

	result, err := diContainer.Get("test")
	if err != nil {
		t.Error(err.Error())
	}

	if result == nil {
		t.Error("Expected object not stored in DI Container")
	}

	if result != expected {
		t.Errorf("Expected object is not same as stored in DI Container")
	}
}

func TestDI_Set_BySameKey(t *testing.T) {
	diContainer.Clear()

	sample := &TestType{}
	if err := diContainer.Set("test", sample); err != nil {
		t.Error(err.Error())
	}

	anotherSample := &AnotherTestType{}
	if err := diContainer.Set("test", anotherSample); err == nil {
		t.Error("Expected to get error while storing by same key to DI Container")
	}
}

func TestDI_MustSet(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Unexpected panic: %v", r)
		}
	}()

	diContainer.Clear()

	expected := &TestType{}
	diContainer.MustSet("test", expected)

	result, err := diContainer.Get("test")
	if err != nil {
		t.Error(err.Error())
	}

	if result == nil {
		t.Error("Expected object not stored in DI Container")
	}

	if result != expected {
		t.Errorf("Expected object is not same as stored in DI Container")
	}
}

func TestDI_MustSet_BySameKey(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered from expected panic: %v", r)
			return
		}

		t.Errorf("Expected panic for storing by same key to DI Container")
	}()

	diContainer.Clear()

	sample := &TestType{}
	diContainer.MustSet("test", sample)

	anotherSample := &AnotherTestType{}
	diContainer.MustSet("test", anotherSample)
}

func TestDI_Get(t *testing.T) {
	diContainer.Clear()

	expected := &TestType{}
	if err := diContainer.Set("test", expected); err != nil {
		t.Error(err.Error())
	}

	result, err := diContainer.Get("test")
	if err != nil {
		t.Error(err.Error())
	}

	if result == nil {
		t.Error("Expected object not stored in DI Container")
	}

	if result != expected {
		t.Errorf("Expected object is not same as stored in DI Container")
	}
}

func TestDI_MustGet(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Unexpected panic: %v", r)
		}
	}()

	diContainer.Clear()

	sample := &TestType{}
	diContainer.MustSet("test", sample)

	result := diContainer.MustGet("test")

	if result == nil {
		t.Error("Expected object not stored in DI Container")
	}

	if result != sample {
		t.Errorf("Expected object is not same as stored in DI Container")
	}
}

func TestDI_Get_Unknown(t *testing.T) {
	diContainer.Clear()

	result, err := diContainer.Get("unknown")
	if err == nil {
		t.Errorf("Expected error for retrieving by unknown key from DI Container")
	}

	if result != nil {
		t.Errorf("Unexpected object by unknown key: %v", reflect.TypeOf(result))
	}
}

func TestDI_MustGet_Unknown(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered from expected panic: %v", r)
			return
		}

		t.Errorf("Expected panic for retrieving by unknown key from DI Container")
	}()

	diContainer.Clear()

	diContainer.MustGet("unknown")
}
