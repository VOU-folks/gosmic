package di

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

type DI struct {
	ctx         context.Context
	containerMu sync.Mutex
	container   map[string]interface{}
}

func NewDI(ctx context.Context) *DI {
	return &DI{
		ctx: ctx,

		containerMu: sync.Mutex{},
		container:   make(map[string]interface{}),
	}
}

func (di *DI) Set(key string, dependency interface{}) error {
	di.containerMu.Lock()
	defer di.containerMu.Unlock()

	if dep := di.container[key]; dep != nil {
		return fmt.Errorf("\"%v\" already injected to DI container and holds \"%v\"", key, reflect.TypeOf(dep))
	}
	di.container[key] = dependency

	return nil
}

func (di *DI) MustSet(key string, dependency interface{}) {
	if err := di.Set(key, dependency); err != nil {
		panic(err.Error())
	}
}

func (di *DI) Get(key string) (interface{}, error) {
	di.containerMu.Lock()
	defer di.containerMu.Unlock()

	var dep interface{}

	if dep = di.container[key]; dep == nil {
		return nil, fmt.Errorf("\"%v\" not injected to DI container", key)
	}
	return dep, nil
}

func (di *DI) MustGet(key string) interface{} {
	var (
		dep interface{}
		err error
	)

	if dep, err = di.Get(key); err != nil {
		panic(err.Error())
	}
	return dep
}

func (di *DI) Clear() {
	di.containerMu.Lock()
	defer di.containerMu.Unlock()

	di.container = make(map[string]interface{})
}

func (di *DI) Size() int {
	return len(di.container)
}

func (di *DI) IsEmpty() bool {
	return len(di.container) == 0
}
