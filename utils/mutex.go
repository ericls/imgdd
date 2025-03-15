package utils

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"hash/fnv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type MutexLock interface {
	AcquireLock() (bool, error)
	ReleaseLock() error
}

type PgAdvisoryLock struct {
	db       *sql.DB
	lockKey  int64
	acquired bool
	mu       sync.Mutex
}

func GenerateLockKey(identifier string) int64 {
	hasher := fnv.New64a()
	hasher.Write([]byte(identifier))
	return int64(hasher.Sum64())
}

func NewPgAdvisoryLock(db *sql.DB, lockKey int64) *PgAdvisoryLock {
	return &PgAdvisoryLock{db: db, lockKey: lockKey}
}

func (p *PgAdvisoryLock) AcquireLock() (bool, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var acquired bool
	err := p.db.QueryRow("SELECT pg_try_advisory_lock($1)", p.lockKey).Scan(&acquired)
	if err != nil {
		return false, err
	}
	p.acquired = acquired
	return acquired, nil
}

func (p *PgAdvisoryLock) ReleaseLock() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.acquired {
		return errors.New("lock was not acquired")
	}

	_, err := p.db.Exec("SELECT pg_advisory_unlock($1)", p.lockKey)
	if err == nil {
		p.acquired = false
	}
	return err
}

type RedisLock struct {
	client   *redis.Client
	key      string
	ttl      time.Duration
	acquired bool
	mu       sync.Mutex
}

func NewRedisLock(client *redis.Client, key string, ttl time.Duration) *RedisLock {
	return &RedisLock{
		client: client,
		key:    key,
		ttl:    ttl,
	}
}

func (r *RedisLock) AcquireLock() (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	ctx := context.Background()

	acquired, err := r.client.SetNX(ctx, r.key, "locked", r.ttl).Result()
	if err != nil {
		return false, err
	}

	r.acquired = acquired
	return acquired, nil
}

func (r *RedisLock) ReleaseLock() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.acquired {
		return errors.New("lock was not acquired")
	}

	ctx := context.Background()
	_, err := r.client.Del(ctx, r.key).Result()
	if err == nil {
		r.acquired = false
	}
	return err
}

func RunWithLock(lock MutexLock, fn func() error) error {
	locked, err := lock.AcquireLock()
	if err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}
	if !locked {
		return errors.New("another instance is already running")
	}
	defer lock.ReleaseLock()

	return fn()
}
