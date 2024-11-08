package fixtures

import (
	"reflect"
	"strings"
)

type Factory[T any] interface {
	New(key string, opts ...func(*T)) *T
}
type Getter[T any] interface {
	Get(key string, opts ...func(*T)) *T
}

type Cache[T any] struct {
	Factory[T]
	instances map[string]*T
}

func NewCache[T any](f Factory[T]) *Cache[T] {
	return &Cache[T]{
		Factory:   f,
		instances: make(map[string]*T),
	}
}

func (c *Cache[T]) Get(key string, opts ...func(*T)) *T {
	r, ok := c.instances[key]
	if ok {
		return r
	}
	r = c.Factory.New(key, opts...)
	c.instances[key] = r
	return r
}

type FactoryDispatcher[T any] map[string]func(...func(*T)) *T

func (d FactoryDispatcher[T]) New(key string, opts ...func(*T)) *T {
	fn := d[key]
	if fn == nil {
		return nil
	}
	return fn(opts...)
}

// NewFactoryDispatcher creates a dispatcher for the factory.
// The factory must have NewXXX methods that have ...func(*T) and return *T.
func NewFactoryDispatcher[T any](factory any) FactoryDispatcher[T] {
	res := map[string]func(...func(*T)) *T{}
	refFactory := reflect.ValueOf(factory)
	refFactoryType := refFactory.Type()
	numMethod := refFactoryType.NumMethod()
	for i := 0; i < numMethod; i++ {
		m := refFactoryType.Method(i)
		if !strings.HasPrefix(m.Name, "New") {
			continue
		}
		if m.Type.Kind() != reflect.Func {
			continue
		}
		numIn := m.Type.NumIn()
		if numIn != 2 { // 1st is receiver, 2nd is variadic
			continue
		}
		inType1 := m.Type.In(1)
		if inType1.Kind() != reflect.Slice {
			continue
		}
		inType1Elem := inType1.Elem()
		if inType1Elem.Kind() != reflect.Func {
			continue
		}
		if inType1Elem.NumIn() != 1 {
			continue
		}
		if inType1Elem.In(0) != reflect.TypeOf((*T)(nil)) {
			continue
		}
		if inType1Elem.NumOut() != 0 {
			continue
		}
		numOut := m.Type.NumOut()
		if numOut != 1 {
			continue
		}
		outType0 := m.Type.Out(0)
		if outType0 != reflect.TypeOf((*T)(nil)) {
			continue
		}
		key := strings.TrimPrefix(m.Name, "New")
		res[key] = refFactory.Method(i).Interface().(func(...func(*T)) *T)
	}
	return res
}

type Fixtures[T any] struct {
	FactoryDispatcher[T]
	*Cache[T]
}

func NewFixtures[T any](impl any) *Fixtures[T] {
	r := &Fixtures[T]{}
	r.FactoryDispatcher = NewFactoryDispatcher[T](impl)
	r.Cache = NewCache[T](r)
	return r
}
