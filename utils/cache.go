package utils

import (
	"container/list"
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

// ValidationCache provides caching of validation results by file hash
type ValidationCache interface {
	Get(fileHash string) (interface{}, bool)
	Set(fileHash string, result interface{}, ttl time.Duration) error
	Clear() error
	Stats() CacheStats
}

// CacheStats provides cache performance statistics
type CacheStats struct {
	Size       int           `json:"size"`
	MaxSize    int           `json:"maxSize"`
	Hits       int64         `json:"hits"`
	Misses     int64         `json:"misses"`
	Evictions  int64         `json:"evictions"`
	HitRate    float64       `json:"hitRate"`
	AverageAge time.Duration `json:"averageAge"`
}

// CachedEntry represents a cached validation result with metadata
type CachedEntry struct {
	Result      interface{}
	CachedAt    time.Time
	ExpiresAt   time.Time
	FileHash    string
	AccessCount int64
	element     *list.Element // For LRU tracking
}

// MemoryValidationCache implements validation caching using in-memory storage with LRU eviction
type MemoryValidationCache struct {
	cache    map[string]*CachedEntry
	lruList  *list.List
	mutex    sync.RWMutex
	maxSize  int
	maxBytes int64

	// Statistics
	hits      int64
	misses    int64
	evictions int64
}

// MemoryCacheOptions configures the memory cache behavior
type MemoryCacheOptions struct {
	MaxEntries int   // Maximum number of cached entries
	MaxBytes   int64 // Maximum memory usage in bytes (approximate)
}

// DefaultMemoryCacheOptions returns sensible defaults for memory caching
func DefaultMemoryCacheOptions() *MemoryCacheOptions {
	return &MemoryCacheOptions{
		MaxEntries: 1000,     // Limit number of cached validations
		MaxBytes:   50 << 20, // 50MB max memory usage
	}
}

// NewMemoryValidationCache creates a new memory-based validation cache
func NewMemoryValidationCache(opts *MemoryCacheOptions) *MemoryValidationCache {
	if opts == nil {
		opts = DefaultMemoryCacheOptions()
	}

	return &MemoryValidationCache{
		cache:    make(map[string]*CachedEntry),
		lruList:  list.New(),
		maxSize:  opts.MaxEntries,
		maxBytes: opts.MaxBytes,
	}
}

// Get retrieves a cached validation result by file hash
func (c *MemoryValidationCache) Get(fileHash string) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.cache[fileHash]
	if !exists {
		c.misses++
		return nil, false
	}

	// Check if entry has expired
	if time.Now().After(entry.ExpiresAt) {
		c.removeEntry(fileHash, entry)
		c.misses++
		return nil, false
	}

	// Move to front of LRU list (mark as recently used)
	c.lruList.MoveToFront(entry.element)
	entry.AccessCount++
	c.hits++

	return entry.Result, true
}

// Set stores a validation result in the cache
func (c *MemoryValidationCache) Set(fileHash string, result interface{}, ttl time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Check if entry already exists
	if existing, exists := c.cache[fileHash]; exists {
		// Update existing entry
		existing.Result = result
		existing.CachedAt = time.Now()
		existing.ExpiresAt = time.Now().Add(ttl)
		c.lruList.MoveToFront(existing.element)
		return nil
	}

	// Create new entry
	entry := &CachedEntry{
		Result:      result,
		CachedAt:    time.Now(),
		ExpiresAt:   time.Now().Add(ttl),
		FileHash:    fileHash,
		AccessCount: 0,
	}

	// Add to LRU list
	entry.element = c.lruList.PushFront(fileHash)
	c.cache[fileHash] = entry

	// Evict if we exceed size limits
	c.evictIfNeeded()

	return nil
}

// Clear removes all cached validation results
func (c *MemoryValidationCache) Clear() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache = make(map[string]*CachedEntry)
	c.lruList = list.New()
	c.hits = 0
	c.misses = 0
	c.evictions = 0

	return nil
}

// Stats returns cache performance statistics
func (c *MemoryValidationCache) Stats() CacheStats {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	totalAccess := c.hits + c.misses
	hitRate := 0.0
	if totalAccess > 0 {
		hitRate = float64(c.hits) / float64(totalAccess)
	}

	// Calculate average age of entries
	var totalAge time.Duration
	entryCount := 0
	now := time.Now()
	for _, entry := range c.cache {
		if !now.After(entry.ExpiresAt) { // Only count non-expired entries
			totalAge += now.Sub(entry.CachedAt)
			entryCount++
		}
	}

	averageAge := time.Duration(0)
	if entryCount > 0 {
		averageAge = totalAge / time.Duration(entryCount)
	}

	return CacheStats{
		Size:       len(c.cache),
		MaxSize:    c.maxSize,
		Hits:       c.hits,
		Misses:     c.misses,
		Evictions:  c.evictions,
		HitRate:    hitRate,
		AverageAge: averageAge,
	}
}

// evictIfNeeded removes old entries when cache limits are exceeded
func (c *MemoryValidationCache) evictIfNeeded() {
	// Remove expired entries first
	c.removeExpiredEntries()

	// Evict LRU entries if we still exceed limits
	for len(c.cache) > c.maxSize {
		if c.lruList.Len() == 0 {
			break
		}

		// Remove least recently used entry
		oldest := c.lruList.Back()
		if oldest != nil {
			fileHash := oldest.Value.(string)
			if entry, exists := c.cache[fileHash]; exists {
				c.removeEntry(fileHash, entry)
				c.evictions++
			}
		}
	}
}

// removeExpiredEntries removes all expired cache entries
func (c *MemoryValidationCache) removeExpiredEntries() {
	now := time.Now()
	toRemove := make([]string, 0)

	for fileHash, entry := range c.cache {
		if now.After(entry.ExpiresAt) {
			toRemove = append(toRemove, fileHash)
		}
	}

	for _, fileHash := range toRemove {
		if entry, exists := c.cache[fileHash]; exists {
			c.removeEntry(fileHash, entry)
		}
	}
}

// removeEntry removes a single entry from cache and LRU list
func (c *MemoryValidationCache) removeEntry(fileHash string, entry *CachedEntry) {
	delete(c.cache, fileHash)
	if entry.element != nil {
		c.lruList.Remove(entry.element)
	}
}

// CalculateFileHash computes SHA256 hash of file content
func CalculateFileHash(content []byte) string {
	hash := sha256.Sum256(content)
	return fmt.Sprintf("%x", hash)
}
