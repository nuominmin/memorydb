package memorydb

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrKeyNotFound       = errors.New("key not found")
	ErrKeyExpired        = errors.New("key expired")
	ErrDBClosed          = errors.New("memory db is closed")
	ErrTypeAssertionFail = errors.New("type assertion failed")
)

type item struct {
	value      interface{}
	expiration int64
}

// MemoryDB 是一个内存数据存储
type MemoryDB struct {
	data   map[string]item
	mutex  sync.RWMutex
	closed bool
	stopCh chan struct{}
	wg     sync.WaitGroup
}

// NewMemoryDB 创建一个新的 MemoryDB
func NewMemoryDB() *MemoryDB {
	db := &MemoryDB{
		data:   make(map[string]item),
		stopCh: make(chan struct{}),
	}
	db.wg.Add(1)
	go db.cleanupExpiredKeys()
	return db
}

// Set 设置一个键值对，可以选择设置过期时间（以秒为单位）
func (db *MemoryDB) Set(key string, value interface{}, ttl int64) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if db.closed {
		return ErrDBClosed
	}

	expiration := int64(0)
	if ttl > 0 {
		expiration = time.Now().UnixNano() + ttl*int64(time.Second)
	}

	db.data[key] = item{
		value:      value,
		expiration: expiration,
	}
	return nil
}

// Get 获取一个键的值，如果键不存在或者已过期则返回错误
func (db *MemoryDB) Get(key string) (interface{}, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	if db.closed {
		return nil, ErrDBClosed
	}

	it, found := db.data[key]
	if !found {
		return nil, ErrKeyNotFound
	}

	if it.expiration > 0 && time.Now().UnixNano() > it.expiration {
		delete(db.data, key)
		return nil, ErrKeyExpired
	}

	return it.value, nil
}

// Del 删除一个键
func (db *MemoryDB) Del(key string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if db.closed {
		return ErrDBClosed
	}

	delete(db.data, key)
	return nil
}

// Expire 设置一个键的过期时间（以秒为单位）
func (db *MemoryDB) Expire(key string, ttl int64) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if db.closed {
		return ErrDBClosed
	}

	it, found := db.data[key]
	if !found {
		return ErrKeyNotFound
	}

	it.expiration = 0
	if ttl > 0 {
		it.expiration = time.Now().UnixNano() + ttl*int64(time.Second)
	}

	db.data[key] = it
	return nil
}

// cleanupExpiredKeys 定期清理过期的键
func (db *MemoryDB) cleanupExpiredKeys() {
	defer db.wg.Done()
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			db.mutex.Lock()
			for key, it := range db.data {
				if it.expiration > 0 && time.Now().UnixNano() > it.expiration {
					delete(db.data, key)
				}
			}
			db.mutex.Unlock()
		case <-db.stopCh:
			return
		}
	}
}

// Close 关闭 MemoryDB 并且清理任务
func (db *MemoryDB) Close() {
	db.mutex.Lock()
	db.closed = true
	close(db.stopCh)
	db.mutex.Unlock()
	db.wg.Wait()

	db.mutex.Lock()
	db.data = make(map[string]item)
	db.mutex.Unlock()
}

// Set 是一个泛型函数，用于设置键值对
func Set[T any](db *MemoryDB, key string, value T, ttl int64) error {
	return db.Set(key, value, ttl)
}

// Get 是一个泛型函数，用于获取键值对
func Get[T any](db *MemoryDB, key string) (T, error) {
	var zeroValue T
	value, err := db.Get(key)
	if err != nil {
		return zeroValue, err
	}
	result, ok := value.(T)
	if !ok {
		return zeroValue, ErrTypeAssertionFail
	}
	return result, nil
}
