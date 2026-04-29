package limiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisLimiter implements distributed rate limiting using Redis + Lua sliding window.
type RedisLimiter struct {
	client *redis.Client
	rate   float64
	burst  int
	prefix string
}

// NewRedisLimiter creates a new distributed Redis rate limiter.
func NewRedisLimiter(client *redis.Client, r float64, burst int) *RedisLimiter {
	return &RedisLimiter{
		client: client,
		rate:   r,
		burst:  burst,
		prefix: "ratelimit:",
	}
}

// slidingWindowScript is a Lua script for sliding window rate limiting.
var slidingWindowScript = redis.NewScript(`
local key = KEYS[1]
local now = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local limit = tonumber(ARGV[3])

-- Remove expired entries
redis.call('ZREMRANGEBYSCORE', key, 0, now - window)

-- Count current entries
local count = redis.call('ZCARD', key)

if count < limit then
    -- Add new entry
    redis.call('ZADD', key, now, now .. '-' .. math.random(1000000))
    redis.call('EXPIRE', key, math.ceil(window / 1000))
    return 1
else
    return 0
end
`)

func (l *RedisLimiter) Allow(key string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	redisKey := l.prefix + key
	now := time.Now().UnixMilli()
	window := int64(l.rate * 1000) // window in ms
	if window <= 0 {
		window = 1000
	}

	result, err := slidingWindowScript.Run(ctx, l.client, []string{redisKey}, now, window, int(l.rate)).Int()
	if err != nil {
		// On Redis error, allow the request (fail open)
		return true
	}

	return result == 1
}
