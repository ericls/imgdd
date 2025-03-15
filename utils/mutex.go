package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"hash/fnv"
	"sync"
)

type DistributedLock interface {
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

func RunWithLock(lock DistributedLock, fn func() error) error {
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
