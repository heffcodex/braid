package braid

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

var _ fiber.Storage = (*StorageRedis)(nil)

type StorageRedis struct {
	rc          redis.UniversalClient
	ctxProvider ContextProvider
	keysPrefix  string
	flushDBMode bool
	allowClose  bool
}

func NewStorageRedis(rc redis.UniversalClient, ctxProvider ContextProvider, keysPrefix string) *StorageRedis {
	return &StorageRedis{rc: rc, ctxProvider: ctxProvider, keysPrefix: keysPrefix}
}

func (s *StorageRedis) EnableFlushDB(v bool) { s.flushDBMode = v }
func (s *StorageRedis) AllowClose(v bool)    { s.allowClose = v }

func (s *StorageRedis) WithPrefix(prefix string) *StorageRedis {
	clone := s.clone()
	clone.keysPrefix += prefix
	return clone
}

func (s *StorageRedis) Get(key string) ([]byte, error) {
	if key == "" {
		return nil, nil
	}

	val, err := s.rc.Get(s.ctxProvider(), s.makeKey(key)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}

	return val, err
}

func (s *StorageRedis) Set(key string, val []byte, exp time.Duration) error {
	if key == "" || len(val) == 0 {
		return nil
	}

	return s.rc.Set(s.ctxProvider(), s.makeKey(key), val, exp).Err()
}

func (s *StorageRedis) Delete(key string) error {
	if key == "" {
		return nil
	}

	return s.rc.Del(s.ctxProvider(), s.makeKey(key)).Err()
}

func (s *StorageRedis) Reset() error {
	ctx := s.ctxProvider()

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
		keys, cursor, err = s.rc.Scan(ctx, cursor, s.makeKey("*"), count).Result()
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

func (s *StorageRedis) Close() error {
	if s.allowClose {
		return s.rc.Close()
	}

	return nil
}

func (s *StorageRedis) makeKey(key string) string {
	return s.keysPrefix + key
}

func (s *StorageRedis) clone() *StorageRedis {
	return &StorageRedis{
		rc:          s.rc,
		ctxProvider: s.ctxProvider,
		keysPrefix:  s.keysPrefix,
		flushDBMode: s.flushDBMode,
		allowClose:  s.allowClose,
	}
}
