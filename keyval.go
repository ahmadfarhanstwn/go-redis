package main

import "sync"

type KeyVal struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func NewKeyVal() *KeyVal {
	return &KeyVal{
		data: map[string][]byte{},
	}
}

func (kv *KeyVal) Set(key string, val []byte) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.data[key] = val
	return nil
}

func (kv *KeyVal) Get(key string) ([]byte, bool) {
	kv.mu.RLock()
	defer kv.mu.Unlock()
	val, ok := kv.data[key]
	return val, ok
}
