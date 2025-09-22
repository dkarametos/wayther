package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// CacheEntry represents a single entry in the cache.
type CacheEntry struct {
	Timestamp time.Time          `json:"timestamp"`
	Weather   *WeatherAPIResponse `json:"weather"`
}

// Cache represents the cache of weather data.
type Cache struct {
	Entries  map[string]CacheEntry `json:"entries"`
	filePath string
}

// NewCache creates a new Cache instance and loads the cache from disk.
func NewCache(configPath string) (*Cache, error) {
	cachePath := filepath.Join(filepath.Dir(configPath), "cache.json")
	cache := &Cache{
		Entries:  make(map[string]CacheEntry),
		filePath: cachePath,
	}
	if err := cache.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return cache, nil
}

// load reads the cache file from disk and unmarshals it into the Cache struct.
func (c *Cache) load() error {
	data, err := os.ReadFile(c.filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &c.Entries)
}

// save writes the cache to disk as a JSON file.
func (c *Cache) save() error {
	data, err := json.MarshalIndent(c.Entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.filePath, data, 0644)
}

// Get retrieves a cache entry for a given location.
func (c *Cache) Get(location string) (*CacheEntry, bool) {
	entry, found := c.Entries[location]
	return &entry, found
}

// Set adds or updates a cache entry and saves the cache to disk.
func (c *Cache) Set(location string, weather *WeatherAPIResponse) error {
	c.Clean(time.Hour)
	c.Entries[location] = CacheEntry{
		Timestamp: time.Now(),
		Weather:   weather,
	}
	return c.save()
}

// IsStale checks if the cache entry is older than the given duration.
func (e *CacheEntry) IsStale(duration time.Duration) bool {
	return time.Since(e.Timestamp) > duration
}

// Clean removes stale entries from the cache and saves the cache to disk.
func (c *Cache) Clean(duration time.Duration) {
	for location, entry := range c.Entries {
		if entry.IsStale(duration) {
			delete(c.Entries, location)
		}
	}
	c.save()
}