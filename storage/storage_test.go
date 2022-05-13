package storage

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCache(t *testing.T) {
	size := 100
	cache := NewMemoryCache(uint32(size))

	t.Run("Put", func(t *testing.T) {
		testEntry := &cacheEntity{
			key:  "put",
			data: []byte("put"),
		}

		err := cache.Put(testEntry.Key(), testEntry.Data())
		assert.NoError(t, err, err)

		t.Run("TheSame", func(t *testing.T) {
			err := cache.Put(testEntry.Key(), testEntry.Data())
			assert.EqualError(t, err, ErrEntryAlreadyExist.Error())
		})
	})

	t.Run("Get", func(t *testing.T) {
		testEntry := &cacheEntity{
			key:  "get",
			data: []byte("get"),
		}

		err := cache.Put(testEntry.Key(), testEntry.Data())
		assert.NoError(t, err, err)

		entry, err := cache.Get(testEntry.Key())
		require.NoError(t, err, err)
		assert.EqualValues(t, testEntry.Data(), entry.Data())

		t.Run("NonExistent", func(t *testing.T) {
			entry, err := cache.Get("fake")
			assert.EqualError(t, err, ErrNotFound.Error())
			assert.Nil(t, entry)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		testEntry := &cacheEntity{
			key:  "delete",
			data: []byte("delete"),
		}

		err := cache.Put(testEntry.Key(), testEntry.Data())
		assert.NoError(t, err, err)

		_, err = cache.Get(testEntry.Key())
		assert.NoError(t, err, err)

		err = cache.Delete(testEntry.Key())
		assert.NoError(t, err, err)

		entry, err := cache.Get(testEntry.Key())
		assert.EqualError(t, err, ErrNotFound.Error())
		assert.Nil(t, entry)
	})

	t.Run("DeleteAll", func(t *testing.T) {
		testEntry := &cacheEntity{
			key:  "deleteAll",
			data: []byte("deleteAll"),
		}

		err := cache.Put(testEntry.Key(), testEntry.Data())
		assert.NoError(t, err, err)

		_, err = cache.Get(testEntry.Key())
		assert.NoError(t, err, err)

		err = cache.DeleteAll()
		assert.NoError(t, err, err)

		entry, err := cache.Get(testEntry.Key())
		assert.EqualError(t, err, ErrNotFound.Error())
		assert.Nil(t, entry)
	})
}

func TestCache_ExceededAmount(t *testing.T) {
	size := 50
	cache := NewMemoryCache(uint32(size))
	testEntries := generateTestData(size)

	for _, te := range testEntries {
		err := cache.Put(te.Key(), te.Data())
		require.NoError(t, err, err)
	}

	testEntry := &cacheEntity{
		key:  "exceed",
		data: []byte("exceed"),
	}
	err := cache.Put(testEntry.Key(), testEntry.Data())
	assert.EqualError(t, err, ErrCacheIsFull.Error())
}

func generateTestData(size int) []*cacheEntity {
	testData := make([]*cacheEntity, size)
	for i := 0; i < size; i++ {
		testData[i] = &cacheEntity{
			key:  fmt.Sprintf("key%d", i),
			data: []byte(fmt.Sprintf("data%d", i)),
		}
	}

	return testData
}
