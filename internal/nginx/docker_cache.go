package nginx

import (
	"sync"
	"time"
)

// Cache for Docker detection to avoid repeated calls
type dockerCache struct {
	mu          sync.RWMutex
	containerID string
	lastCheck   time.Time
	cacheTTL    time.Duration
}

var cache = &dockerCache{
	cacheTTL: 5 * time.Second, // Cache for 5 seconds
}

// getCachedContainerID returns cached container ID or detects it
func getCachedContainerID() (string, error) {
	cache.mu.RLock()
	if time.Since(cache.lastCheck) < cache.cacheTTL && cache.containerID != "" {
		containerID := cache.containerID
		cache.mu.RUnlock()
		return containerID, nil
	}
	cache.mu.RUnlock()

	// Cache expired or empty, detect again
	containerID, err := DetectDockerNginx()
	if err != nil {
		return "", err
	}

	// Update cache
	cache.mu.Lock()
	cache.containerID = containerID
	cache.lastCheck = time.Now()
	cache.mu.Unlock()

	return containerID, nil
}

// invalidateCache clears the cache
func invalidateCache() {
	cache.mu.Lock()
	cache.containerID = ""
	cache.lastCheck = time.Time{}
	cache.mu.Unlock()
}
