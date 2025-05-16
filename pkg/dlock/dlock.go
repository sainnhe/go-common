//go:generate mockgen -write_package_comment=false -source=dlock.go -destination=dlock_mock.go -package dlock

// Package dlock implements distributed lock.
package dlock

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/rueidis"
	"github.com/sainnhe/go-common/pkg/constant"
)

// ErrKeyNotExists indicates that the key doesn't exist.
var ErrKeyNotExists = errors.New("key doesn't exist")

// Service is the distributed lock service.
type Service interface {
	// TryAcquire tries to acquire a key without waiting and returns whether it can be acquired.
	TryAcquire(ctx context.Context, key string) (bool, error)

	// Acquire acquires a key.
	// If it's already acquired by others, wait and retry until ctx is cancelled.
	Acquire(ctx context.Context, key string) error

	// Release releases a key.
	// [ErrKeyNotExists] might be returned if it doesn't exist.
	Release(ctx context.Context, key string) error
}

type serviceImpl struct {
	cfg *Config
	rc  rueidis.Client
}

// NewService initializes a new dlock service.
func NewService(cfg *Config, rc rueidis.Client) (Service, error) {
	if cfg == nil || rc == nil {
		return nil, constant.ErrNilDeps
	}
	return &serviceImpl{
		cfg,
		rc,
	}, nil
}

func (s *serviceImpl) TryAcquire(ctx context.Context, key string) (bool, error) {
	err := s.rc.Do(ctx, s.rc.B().Get().Key(s.getKey(key)).Build()).Error()
	switch err {
	case rueidis.Nil:
		return true, nil
	case nil:
		return false, nil
	default:
		return false, err
	}
}

func (s *serviceImpl) Acquire(ctx context.Context, key string) error {
	for {
		err := s.rc.Do(ctx, s.rc.B().
			Set().
			Key(s.getKey(key)).
			Value("1").
			Nx().
			PxMilliseconds(s.cfg.ExpireMs).
			Build()).Error()
		switch err {
		case rueidis.Nil:
			time.Sleep(time.Duration(s.cfg.RetryAfterMs) * time.Millisecond)
			continue
		case nil:
			return nil
		default:
			return err
		}
	}
}

func (s *serviceImpl) Release(ctx context.Context, key string) error {
	v, err := s.rc.Do(ctx, s.rc.B().
		Del().
		Key(s.getKey(key)).
		Build()).AsInt64()
	if err != nil {
		return err
	}
	if v == 1 {
		return nil
	}
	return ErrKeyNotExists
}

func (s *serviceImpl) getKey(key string) string {
	return fmt.Sprintf("%s:%s", s.cfg.Prefix, key)
}
