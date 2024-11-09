package fixtures

import (
	"reflect"
	"strings"
)

// Factory is an interface that defines a method to create a new instance of type T.
type Factory[T any] interface {
	New(key string, opts ...func(*T)) *T
}

// Getter is an interface that defines a method to get an instance of type T.
type Getter[T any] interface {
	Get(key string, opts ...func(*T)) *T
}

// Cache is a struct that implements caching for instances of type T.
type Cache[T any] struct {
	Factory[T]
	instances map[string]*T
}

// NewCache creates a new Cache with the given Factory.
func NewCache[T any](f Factory[T]) *Cache[T] {
	return &Cache[T]{
		Factory:   f,
		instances: make(map[string]*T),
	}
}

// Get retrieves an instance from the cache or creates a new one using the Factory.
func (c *Cache[T]) Get(key string, opts ...func(*T)) *T {
	r, ok := c.instances[key]
	if ok {
		return r
	}
	r = c.Factory.New(key, opts...)
	c.instances[key] = r
	return r
}

// FactoryDispatcher is a map that dispatches factory methods based on a key.
type FactoryDispatcher[T any] map[string]func(...func(*T)) *T

// New creates a new instance of type T using the factory method associated with the key.
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

// Fixtures is a struct that combines a FactoryDispatcher and a Cache.
type Fixtures[T any] struct {
	FactoryDispatcher[T]
	*Cache[T]
}

// NewFixtures creates a new Fixtures instance with the given implementation.
func NewFixtures[T any](impl any) *Fixtures[T] {
	r := &Fixtures[T]{}
	r.FactoryDispatcher = NewFactoryDispatcher[T](impl)
	r.Cache = NewCache[T](r)
	return r
}
