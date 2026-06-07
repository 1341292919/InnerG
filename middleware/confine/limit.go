package confine

import (
	"InnerG/dao/cache"
	"InnerG/pack"
	"InnerG/pkg/errno"
	"InnerG/pkg/logger"
	"bytes"
	"context"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	c *redis.Client
}

type Rule struct {
	Name   string
	Window time.Duration
	Max    int64
	Key    KeyFunc
}

type KeyFunc func(c *gin.Context) string

func NewLimiter(client *redis.Client) *Limiter {
	return &Limiter{c: client}
}

func Limit(rules ...Rule) gin.HandlerFunc {
	limiter := NewLimiter(cache.GetRedisClient())
	return func(c *gin.Context) {
		for _, rule := range rules {
			keyPart := rule.Key(c)
			if keyPart == "" {
				continue
			}

			key := fmt.Sprintf("rate:%s:%s", rule.Name, keyPart)
			allowed, err := limiter.Allow(c.Request.Context(), key, rule.Window, rule.Max)
			if err != nil {
				if logger.Log != nil {
					logger.Log.Errorf("rate limit check failed: %v", err)
				}
				c.Next()
				return
			}

			if !allowed {
				if logger.Log != nil {
					logger.Log.Errorf("rate limit exceeded: rule=%s key=%s ip=%s path=%s", rule.Name, key, c.ClientIP(), c.Request.URL.Path)
				}
				pack.RespError(c, errno.RateLimitExceeded)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

func (l *Limiter) Allow(ctx context.Context, key string, window time.Duration, max int64) (bool, error) {
	if l.c == nil || max <= 0 || window <= 0 {
		return true, nil
	}

	count, err := l.c.Incr(ctx, key).Result()
	if err != nil {
		return true, err
	}
	if count == 1 {
		if err := l.c.Expire(ctx, key, window).Err(); err != nil {
			return true, err
		}
	}

	return count <= max, nil
}

func ByIP() KeyFunc {
	return func(c *gin.Context) string {
		return c.ClientIP()
	}
}

func ByField(field string) KeyFunc {
	return func(c *gin.Context) string {
		value := strings.TrimSpace(c.PostForm(field))
		if value == "" {
			value = strings.TrimSpace(c.Query(field))
		}
		return value
	}
}

func ByBind[T any](pick func(*T) string) KeyFunc {
	return func(c *gin.Context) string {
		req := bindRequest[T](c)
		if req == nil {
			return ""
		}
		return strings.TrimSpace(pick(req))
	}
}

func bindRequest[T any](c *gin.Context) *T {
	cacheKey := fmt.Sprintf("confine:bind:%s", reflect.TypeOf((*T)(nil)).Elem().String())
	if v, ok := c.Get(cacheKey); ok {
		if req, ok := v.(*T); ok {
			return req
		}
	}

	req := new(T)
	var body []byte
	if c.Request != nil && c.Request.Body != nil {
		var err error
		body, err = io.ReadAll(c.Request.Body)
		if err != nil {
			return nil
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(body))
	}

	_ = c.ShouldBind(req)
	if body != nil {
		c.Request.Body = io.NopCloser(bytes.NewReader(body))
	}

	c.Set(cacheKey, req)
	return req
}

func Compose(keys ...KeyFunc) KeyFunc {
	return func(c *gin.Context) string {
		parts := make([]string, 0, len(keys))
		for _, key := range keys {
			part := key(c)
			if part == "" {
				return ""
			}
			parts = append(parts, part)
		}
		return strings.Join(parts, ":")
	}
}
