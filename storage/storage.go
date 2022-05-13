package storage

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"sync"
)

const DefaultSize = 1000

var (
	ErrNotFound          = errors.New("entity not found")
	ErrCacheIsFull       = errors.New("cache is full")
	ErrEntryAlreadyExist = errors.New("entry already exist")
)

type (
	Entity interface {
		Key() string
		Data() []byte
	}

	Storage interface {
		// Put puts new entity
		Put(key string, data []byte) error
		// Get gets an entity by id
		Get(string) (Entity, error)
		// Delete deletes an entity by id
		Delete(string) error
		// DeleteAll deletes all entities
		DeleteAll() error
	}

	MemoryStorage interface {
		Storage
		// SaveToFile gets file path nad saves to it
		SaveToFile(string) error
		// LoadFromFile gets file path nad loads form file
		LoadFromFile(string) error
	}

	cacheEntity struct {
		key  string
		data []byte
	}

	cache struct {
		cacheLock sync.Mutex
		cache     map[string][]byte
		size      uint32
	}
)

func NewMemoryCache(size uint32) MemoryStorage {
	if size == 0 {
		size = DefaultSize
	}

	return &cache{
		cacheLock: sync.Mutex{},
		cache:     map[string][]byte{},
		size:      size,
	}
}

func (m *cacheEntity) Key() string {
	return m.key
}

func (m *cacheEntity) Data() []byte {
	return m.data
}

func (m *cache) Put(key string, data []byte) error {
	m.cacheLock.Lock()
	defer m.cacheLock.Unlock()

	if uint32(len(m.cache)) == m.size {
		return ErrCacheIsFull
	}

	if _, ok := m.cache[key]; ok {
		return ErrEntryAlreadyExist
	}

	m.cache[key] = data
	return nil
}

func (m *cache) Get(key string) (Entity, error) {
	m.cacheLock.Lock()
	defer m.cacheLock.Unlock()

	data := m.cache[key]
	if data == nil {
		return nil, ErrNotFound
	}

	return &cacheEntity{
		key:  key,
		data: data,
	}, nil
}

func (m *cache) Delete(key string) error {
	m.cacheLock.Lock()
	defer m.cacheLock.Unlock()
	delete(m.cache, key)

	return nil
}

func (m *cache) DeleteAll() error {
	m.cacheLock.Lock()
	defer m.cacheLock.Unlock()
	m.cache = nil

	return nil
}

func (m *cache) SaveToFile(path string) error {
	m.cacheLock.Lock()
	defer m.cacheLock.Unlock()

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer f.Close()

	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(&m.cache)
	if err != nil {
		return err
	}

	_, err = f.Write(b.Bytes())
	return err
}

func (m *cache) LoadFromFile(path string) error {
	m.cacheLock.Lock()
	defer m.cacheLock.Unlock()

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	b := bytes.NewReader(data)
	return json.NewDecoder(b).Decode(&m.cache)
}
