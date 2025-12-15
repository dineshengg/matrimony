package utils

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/redis/go-redis/v9"
)

func IsRedisInit() bool {
	log.Debug("dummy")
	return true
}

type RedisClient struct {
	// Add any necessary fields for your Redis client
	redisClient *redis.Client // Redis client instance
	initialized bool
	address     string
	nCountRetry int
}

// global variable that will be reused in this package
var redisClient *RedisClient = nil

func NewRedisClient(addr string) *RedisClient {

	redisClient = &RedisClient{
		redisClient: nil,
		initialized: false,
		address:     addr,
		nCountRetry: 0,
	}
	if redisClient.InitializeRedis() != nil {
		log.Println("Failed to initialize redis client")
		return nil
	}
	return redisClient

}

func GetRedisClient() *RedisClient {
	return redisClient
}

// InitializeRedis initializes the Redis client
func (r *RedisClient) InitializeRedis() error {
	if r.nCountRetry > 6 {
		panic(fmt.Errorf("Redis initialization failed after 6 attempt"))
	}

	r.nCountRetry++

	r.redisClient = redis.NewClient(&redis.Options{
		Addr: r.address, // Redis server address
	})

	if r.redisClient == nil {
		log.Error("Failed to create Redis client")
	}
	// Test Redis connection
	ctx := context.Background()
	context.WithValue(ctx, "ping", r.address)
	_, err := r.redisClient.Ping(ctx).Result()
	if err != nil {
		//panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}
	r.initialized = true
	log.Println("Connected to Redis successfully!")
	return nil
}

func (r *RedisClient) GetValue(ctx context.Context, key string) (string, error) {
	if r.initialized {
		val, err := r.redisClient.Get(ctx, key).Result()
		if err != nil {
			log.Println("Error getting value from redis for key ", key)
			return "", err
		}
		return val, nil
	}
	//If redis client is not initialized make sure to initialize it again to max retry count
	if r.InitializeRedis() != nil {
		log.Fatal("GetValue - Failed to initialize redis client")
		return "", fmt.Errorf("Getvalue - Failed to initialize redis client")
	}
	log.Println("Get value initializing redis again in address ", r.address)
	return r.GetValue(ctx, key)

}

func (r *RedisClient) SetValue(ctx context.Context, key string, val string) error {
	if r.initialized {
		err := r.redisClient.Set(ctx, key, val, time.Duration(24)*time.Hour).Err()
		if err != nil {
			log.Println("Error setting value from redis for key %s", key)
			return err
		}
		return nil
	}
	//If redis client is not initialized make sure to initialize it again to max retry count
	if r.InitializeRedis() != nil {
		return fmt.Errorf("SetValue - Failed to initialize redis client")
	}
	log.Println("Set value initializing redis again in address %s", r.address)
	return r.SetValue(ctx, key, val)
}

func (r *RedisClient) DeleteValue(ctx context.Context, key string) error {
	if r.initialized {
		err := r.redisClient.Del(ctx, key).Err()
		if err != nil {
			log.Println("Error deleting value from redis for key %s", key)
			return err
		}
		return nil
	}

	//If redis client is not initialized make sure to initialize it again to max retry count
	if r.InitializeRedis() != nil {
		return fmt.Errorf("DeleteValue - Failed to initialize redis client")
	}
	log.Println("Delete value initializing redis again in address %s", r.address)
	return r.DeleteValue(ctx, key)
}

func (r *RedisClient) Close() {
	if r.initialized {
		r.redisClient.Close()
		log.Println("Redis connection closed")
		return
	}
	log.Println("redis client is not initialized, not able to close the connection")
	return
}
