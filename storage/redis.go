package storage

import (
	"time"

	"github.com/redis/go-redis/v9"
)

var _ Storage = (*Redis)(nil)

type RedisOption func(*Redis)

func RedisAllowClose(v ...bool) RedisOption {
	allow := len(v) == 0 || v[0]
	return func(s *Redis) { s.allowClose = allow }
}

func RedisContextProvider(ctxProvider ContextProvider) RedisOption {
	if ctxProvider == nil {
		ctxProvider = ctxProviderDefault
	}

	return func(s *Redis) { s.ctxp = ctxProvider }
}

func RedisEnableFlushDB(v ...bool) RedisOption {
	enable := len(v) == 0 || v[0]
	return func(s *Redis) { s.flushDBMode = enable }
}

func RedisNamespace(ns string) RedisOption {
	return func(s *Redis) { s.ns += ns + ":" }
}

type Redis struct {
	ctxp        ContextProvider
	rc          redis.UniversalClient
	ns          string
	flushDBMode bool
	allowClose  bool
}

func NewRedis(rc redis.UniversalClient, opts ...RedisOption) *Redis {
	s := &Redis{ctxp: ctxProviderDefault, rc: rc}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Redis) Get(key string) ([]byte, error) {
	if key == "" {
		return nil, nil
	}

	val, err := s.rc.Get(s.ctxp(), s.key(key)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}

	return val, err
}

func (s *Redis) Set(key string, val []byte, exp time.Duration) error {
	if key == "" || len(val) == 0 {
		return nil
	}

	return s.rc.Set(s.ctxp(), s.key(key), val, exp).Err()
}

func (s *Redis) Delete(key string) error {
	if key == "" {
		return nil
	}

	return s.rc.Del(s.ctxp(), s.key(key)).Err()
}

func (s *Redis) Reset() error {
	ctx := s.ctxp()

	if s.flushDBMode {
		return s.rc.FlushDB(ctx).Err()
	}

	const count = int64(100)

	allKeys := make([]string, 0, count)
	cursor := uint64(0)

	var (
		keys []string
		err  error
	)

	for {
		keys, cursor, err = s.rc.Scan(ctx, cursor, s.key("*"), count).Result()
		if err != nil {
			return err
		}

		allKeys = append(allKeys, keys...)

		if cursor == 0 {
			break
		}
	}

	if len(allKeys) > 0 {
		return s.rc.Del(ctx, allKeys...).Err()
	}

	return nil
}

func (s *Redis) Close() error {
	if s.allowClose {
		return s.rc.Close()
	}

	return nil
}

func (s *Redis) key(key string) string {
	return s.ns + key
}
