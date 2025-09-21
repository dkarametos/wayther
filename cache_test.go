package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "cache-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	cachePath := filepath.Join(tempDir, "cache.json")

	t.Run("NewCache", func(t *testing.T) {
		// Test creating a new cache when the file doesn't exist
		cache, err := NewCache(cachePath)
		assert.NoError(t, err)
		assert.NotNil(t, cache)
		assert.Empty(t, cache.Entries)
	})

	t.Run("Set and Get", func(t *testing.T) {
		cache, err := NewCache(cachePath)
		assert.NoError(t, err)

		// Create a mock weather response
		mockWeather := &WeatherAPIResponse{
			Location: Location{
				Name: "Test Location",
			},
		}

		// Set the cache entry
		err = cache.Set("Test Location", mockWeather)
		assert.NoError(t, err)

		// Get the cache entry
		entry, found := cache.Get("Test Location")
		assert.True(t, found)
		assert.NotNil(t, entry)
		assert.Equal(t, "Test Location", entry.Weather.Location.Name)
	})

	t.Run("IsStale", func(t *testing.T) {
		cache, err := NewCache(cachePath)
		assert.NoError(t, err)

		// Create a mock weather response
		mockWeather := &WeatherAPIResponse{
			Location: Location{
				Name: "Test Location",
			},
		}

		// Set the cache entry
		err = cache.Set("Test Location", mockWeather)
		assert.NoError(t, err)

		// Get the cache entry
		entry, found := cache.Get("Test Location")
		assert.True(t, found)

		// Check if it's stale (should not be)
		assert.False(t, entry.IsStale(time.Hour))

		// Manipulate the timestamp to make it stale
		entry.Timestamp = time.Now().Add(-2 * time.Hour)
		assert.True(t, entry.IsStale(time.Hour))
	})

	t.Run("Load and Save", func(t *testing.T) {
		cache, err := NewCache(cachePath)
		assert.NoError(t, err)

		// Create a mock weather response
		mockWeather := &WeatherAPIResponse{
			Location: Location{
				Name: "Test Location",
			},
		}

		// Set the cache entry
		err = cache.Set("Test Location", mockWeather)
		assert.NoError(t, err)

		// Create a new cache instance to load from the file
		newCache, err := NewCache(cachePath)
		assert.NoError(t, err)

		// Get the cache entry from the new cache
		entry, found := newCache.Get("Test Location")
		assert.True(t, found)
		assert.NotNil(t, entry)
		assert.Equal(t, "Test Location", entry.Weather.Location.Name)
	})
}
