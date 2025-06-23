package utils

import (
	"sync"

	"github.com/bradfitz/gomemcache/memcache"
	flag "github.com/spf13/viper"
)

var (
	memcacheClient *memcache.Client
	once           sync.Once
	memcacheErr    error
)

// InitMemcache initializes the singleton memcache client using Viper config.
// Call this once at app startup, or it will auto-init on first use with default localhost:11211.
func InitMemcache() {
	once.Do(func() {
		servers := flag.GetStringSlice("memcache.servers")
		if len(servers) == 0 {
			servers = []string{"127.0.0.1:11211"}
		}
		memcacheClient = memcache.New(servers...)

		// Configure timeouts and idle connections if set in config
		if flag.IsSet("memcache.timeout") {
			timeout := flag.GetDuration("memcache.timeout")
			memcacheClient.Timeout = timeout
		}
		if flag.IsSet("memcache.max_idle_conns") {
			memcacheClient.MaxIdleConns = flag.GetInt("memcache.max_idle_conns")
		}
		// Optionally, ping to check connection
		memcacheErr = memcacheClient.Set(&memcache.Item{Key: "ping", Value: []byte("pong")})
		if memcacheErr != nil {
			memcacheClient = nil                // Reset client on error
			memcacheErr = memcache.ErrNoServers // Set a more specific error
		} else {
			// Clean up the test item after ping
			memcacheClient.Delete("ping")
		}
	})
}

// IsMemcacheInit returns error if memcache client is not initialized.
func IsMemcacheInit() error {
	if memcacheClient == nil {
		return memcache.ErrNoServers
	}
	return memcacheErr
}

// GetMemcacheClient returns the singleton memcache client, initializing if needed.
func GetMemcacheClient() *memcache.Client {
	if memcacheClient == nil {
		InitMemcache()
	}
	return memcacheClient
}

// Set sets a key-value pair in memcache.
func MemcacheSet(key string, value []byte, expiration int32) error {
	client := GetMemcacheClient()
	return client.Set(&memcache.Item{Key: key, Value: value, Expiration: expiration})
}

// Get retrieves a value by key from memcache.
func MemcacheGet(key string) ([]byte, error) {
	client := GetMemcacheClient()
	item, err := client.Get(key)
	if err != nil {
		return nil, err
	}
	return item.Value, nil
}

// MultiGet retrieves multiple values by keys from memcache.
func MemcacheMultiGet(keys []string) (map[string][]byte, error) {
	client := GetMemcacheClient()
	items, err := client.GetMulti(keys)
	if err != nil {
		return nil, err
	}
	result := make(map[string][]byte, len(items))
	for k, item := range items {
		result[k] = item.Value
	}
	return result, nil
}
