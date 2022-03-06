package main

import (
	"sync"
)

type KVStorage interface {
	Get(key string) (interface{}, error)
	Put(key string, val interface{}) error
	Delete(key string) error
}


type SafeMap struct {
	mu sync.RWMutex
	v  map[string]interface{}
}


func (m *SafeMap) Get(key string) (interface{}, error) {
	var err error

	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.v[key], err
}

func (m *SafeMap) Put(key string, val interface{}) error {
	var err error
	if m.v == nil {
		m.v = make(map[string]interface{})
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.v[key] = val
	return err
}

func (m *SafeMap) Delete(key string) error {
	var err error
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.v[key]; ok {
		delete(m.v, key)
	}

	return err
}