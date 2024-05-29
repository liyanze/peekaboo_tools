package safemap

import (
	"sync"
)

type Map[K comparable, V any] struct {
	m sync.Map
}

func (s *Map[K, V]) Set(key K, value V) {
	s.m.Store(key, value)

}

func (s *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	val, loaded := s.m.LoadAndDelete(key)
	return convert[V](val, loaded)
}

func (s *Map[K, V]) Get(key K) (V, bool) {
	val, ok := s.m.Load(key)
	return convert[V](val, ok)
}

func (s *Map[K, V]) Range(f func(key K, value V) bool) {
	s.m.Range(func(key, value any) bool {
		return f(key.(K), value.(V))
	})
}

func (s *Map[K, V]) Delete(key K) {
	s.m.Delete(key)
}

func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{}
}

func convert[T any](val any, ok bool) (T, bool) {
	if !ok {
		var temp T
		return temp, ok
	}
	return val.(T), ok
}
