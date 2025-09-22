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
		os.Remove(cachePath) // Ensure clean slate
		// Test creating a new cache when the file doesn't exist
		cache, err := NewCache(cachePath)
		assert.NoError(t, err)
		assert.NotNil(t, cache)
		assert.Empty(t, cache.Entries)
	})

	t.Run("NewCache with corrupted file", func(t *testing.T) {
		os.Remove(cachePath) // Ensure clean slate
		// Create a corrupted cache file
		err = os.WriteFile(cachePath, []byte("this is not valid json"), 0644)
		assert.NoError(t, err)

		cache, err := NewCache(cachePath)
		assert.Error(t, err)
		assert.Nil(t, cache)
		assert.Contains(t, err.Error(), "invalid character")
	})

	t.Run("NewCache with unreadable file", func(t *testing.T) {
		os.Remove(cachePath) // Ensure clean slate
		// Create an empty file to set permissions on
		err = os.WriteFile(cachePath, []byte(""), 0644)
		assert.NoError(t, err)

		// Make the cache file unreadable
		err = os.Chmod(cachePath, 0000) // No read permissions
		assert.NoError(t, err)

		cache, err := NewCache(cachePath)
		assert.Error(t, err)
		assert.Nil(t, cache)
		assert.Contains(t, err.Error(), "permission denied")

		// Restore permissions for cleanup
		os.Chmod(cachePath, 0644)
	})

	t.Run("Set and Get", func(t *testing.T) {
		os.Remove(cachePath) // Ensure clean slate
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
		os.Remove(cachePath) // Ensure clean slate
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
		os.Remove(cachePath) // Ensure clean slate
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

	t.Run("Clean", func(t *testing.T) {
		os.Remove(cachePath) // Ensure clean slate
		cache, err := NewCache(cachePath)
		assert.NoError(t, err)

		mockWeather := &WeatherAPIResponse{
			Location: Location{
				Name: "Stale Location",
			},
		}
		// Add a stale entry
		cache.Entries["Stale Location"] = CacheEntry{
			Timestamp: time.Now().Add(-2 * time.Hour),
			Weather:   mockWeather,
		}

		mockWeather2 := &WeatherAPIResponse{
			Location: Location{
				Name: "Fresh Location",
			},
		}
		// Add a fresh entry
		cache.Entries["Fresh Location"] = CacheEntry{
			Timestamp: time.Now().Add(-30 * time.Minute),
			Weather:   mockWeather2,
		}

		cache.Clean(time.Hour)

		// Stale entry should be removed
		_, foundStale := cache.Get("Stale Location")
		assert.False(t, foundStale)

		// Fresh entry should remain
		_, foundFresh := cache.Get("Fresh Location")
		assert.True(t, foundFresh)

		// Test cleaning an empty cache
		os.Remove(cachePath)
		emptyCache, err := NewCache(cachePath)
		assert.NoError(t, err)
		emptyCache.Clean(time.Hour)
		assert.Empty(t, emptyCache.Entries)

		// Test cleaning with all entries stale
		os.Remove(cachePath)
		allStaleCache, err := NewCache(cachePath)
		assert.NoError(t, err)
		allStaleCache.Entries["Stale1"] = CacheEntry{Timestamp: time.Now().Add(-2 * time.Hour), Weather: mockWeather}
		allStaleCache.Entries["Stale2"] = CacheEntry{Timestamp: time.Now().Add(-3 * time.Hour), Weather: mockWeather}
		allStaleCache.Clean(time.Hour)
		assert.Empty(t, allStaleCache.Entries)

		// Test cleaning with no entries stale
		os.Remove(cachePath)
		noStaleCache, err := NewCache(cachePath)
		assert.NoError(t, err)
		noStaleCache.Entries["Fresh1"] = CacheEntry{Timestamp: time.Now().Add(-30 * time.Minute), Weather: mockWeather}
		noStaleCache.Entries["Fresh2"] = CacheEntry{Timestamp: time.Now().Add(-45 * time.Minute), Weather: mockWeather}
		noStaleCache.Clean(time.Hour)
		assert.Len(t, noStaleCache.Entries, 2)
	})

	t.Run("Set error handling", func(t *testing.T) {
		os.Remove(cachePath) // Ensure clean slate
		cache, err := NewCache(cachePath)
		assert.NoError(t, err)

		mockWeather := &WeatherAPIResponse{
			Location: Location{
				Name: "Initial Location",
			},
		}
		// Set an initial entry to create the cache file
		err = cache.Set("Initial Location", mockWeather)
		assert.NoError(t, err)

		mockWeatherError := &WeatherAPIResponse{
			Location: Location{
				Name: "Error Location",
			},
		}

		// Make the cache file unwriteable to simulate an error during save
		err = os.Chmod(cachePath, 0444) // Read-only permissions
		assert.NoError(t, err)

		// Attempt to set another entry, which should now fail
		err = cache.Set("Error Location", mockWeatherError)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "permission denied")

		// Restore permissions for other tests
		os.Chmod(cachePath, 0644)
	})

	t.Run("Save error handling", func(t *testing.T) {
		os.Remove(cachePath) // Ensure clean slate
		cache, err := NewCache(cachePath)
		assert.NoError(t, err)

		mockWeather := &WeatherAPIResponse{
			Location: Location{
				Name: "Save Error Location",
			},
		}
		// Add an entry to the cache
		cache.Entries["Save Error Location"] = CacheEntry{
			Timestamp: time.Now(),
			Weather:   mockWeather,
		}
		// Save once to create the file
		err = cache.save()
		assert.NoError(t, err)

		// Make the cache file unwriteable to simulate an error during save
		err = os.Chmod(cachePath, 0444) // Read-only permissions
		assert.NoError(t, err)

		// Directly call save, which should now fail
		err = cache.save()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "permission denied")

		// Restore permissions for cleanup
		os.Chmod(cachePath, 0644)
	})
}
