package llm

import (
	"context"
	"daily-dashboard-backend/src/data"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	RedisClient *redis.Client
}

// Establish the initial Redis Client Connection
func CreateRedisClient() (*RedisClient, error) {
	// Retrieve Environment Variables
	redisUri := os.Getenv("REDIS_URI")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDb, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve RedisDB: %w", err)
	}

	// Create & Verify Redis Client
	client := redis.NewClient(&redis.Options{
		Addr:     redisUri,
		Password: redisPassword,
		DB:       redisDb,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("unable to connect to Redis: %w", err)
	}

	// Verified - Return our Redis Client
	return &RedisClient{
		RedisClient: client,
	}, nil
}

// Proper cleanup of the Redis Cache
func (r *RedisClient) Terminate() error {
	return r.RedisClient.FlushDB(context.Background()).Err()
}

// convo *[]data.Message
func (r *RedisClient) SetConversationData(idHex string, convo *data.Conversation) error {
	// Convert our Conversation to a Byte Array
	convoBytes, err := json.Marshal(convo)
	if err != nil {
		return fmt.Errorf("unable to marshall conversation data for Redis: %w", err)
	}

	// Write it to Redis
	return r.RedisClient.Set(context.Background(), idHex, convoBytes, 0).Err()
}

func (r *RedisClient) GetConversationData(idHex string) (*data.Conversation, error) {
	// Retrieve the Conversation Bytes from Redis
	convoBytes, err := r.RedisClient.Get(context.Background(), idHex).Bytes()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve conversation data from Redis: %w", err)
	}

	// Convert the Byte Array into our Conversation Array
	var convo data.Conversation
	err = json.Unmarshal(convoBytes, &convo)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshall convoBytes from Redis: %w", err)
	}
	return &convo, nil
}

func (r *RedisClient) RemoveConversationData(idHex string) error {
	exists, err := r.RedisClient.Exists(context.Background(), idHex).Result()
	if err != nil {
		return fmt.Errorf("unable to verify whether the ObjectId exists in Redis: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("conversation has not been loaded into Redis, unable to remove it")
	}

	if err := r.RedisClient.Del(context.Background(), idHex).Err(); err != nil {
		return fmt.Errorf("unable to delete the Conversation from Redis: %w", err)
	}
	return nil
}
