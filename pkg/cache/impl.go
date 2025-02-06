package cache

import (
	"context"
	"fmt"

	"github.com/teamsorghum/go-common/pkg/constant"
	"github.com/teamsorghum/go-common/pkg/log"
	"github.com/valkey-io/valkey-go"
)

type proxyImpl struct {
	c valkey.Client
	l log.Logger
}

const (
	attrKey   = "key"
	attrValue = "value"
)

// NewProxyImpl initializes a new cache proxy.
func NewProxyImpl(cfg *Config, logger log.Logger) (proxy Proxy, cleanup func(), err error) {
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)},
		Username:    cfg.Username,
		Password:    cfg.Password,
	})
	if err != nil {
		return
	}
	proxy = &proxyImpl{
		l: logger.WithAttrs(constant.LogAttrAPI, "cache"),
		c: client,
	}
	cleanup = func() {
		if client != nil {
			client.Close()
		}
	}
	return
}

func (p *proxyImpl) Set(ctx context.Context, key, value string) error {
	err := p.c.Do(ctx, p.c.B().Set().Key(key).Value(value).Build()).Error()
	p.l.WithContext(ctx).Debug("Execute set command.", attrKey, key, attrValue, value, constant.LogAttrError, err)
	return err
}

func (p *proxyImpl) Setex(ctx context.Context, key, value string, seconds int64) error {
	err := p.c.Do(ctx, p.c.B().Setex().Key(key).Seconds(seconds).Value(value).Build()).Error()
	p.l.WithContext(ctx).Debug("Execute setex command.",
		attrKey, key, attrValue, value, "seconds", seconds, constant.LogAttrError, err)
	return err
}

func (p *proxyImpl) Get(ctx context.Context, key string) valkey.ValkeyResult {
	r := p.c.Do(ctx, p.c.B().Get().Key(key).Build())
	p.l.WithContext(ctx).Debug("Execute get command.", attrKey, key, "result", r.String())
	return r
}

func (p *proxyImpl) Delete(ctx context.Context, key string) error {
	err := p.c.Do(ctx, p.c.B().Del().Key(key).Build()).Error()
	p.l.WithContext(ctx).Debug("Execute del command.", attrKey, key, constant.LogAttrError, err)
	return err
}

func (p *proxyImpl) Expire(ctx context.Context, key string, seconds int64) error {
	err := p.c.Do(ctx, p.c.B().Expire().Key(key).Seconds(seconds).Build()).Error()
	p.l.WithContext(ctx).Debug("Execute expire command.", attrKey, key, "seconds", seconds, constant.LogAttrError, err)
	return err
}

func (p *proxyImpl) Incr(ctx context.Context, key string) (int64, error) {
	r := p.c.Do(ctx, p.c.B().Incr().Key(key).Build())
	p.l.WithContext(ctx).Debug("Execute incr command.", attrKey, key, "result", r.String())
	return r.AsInt64()
}

func (p *proxyImpl) IncrBy(ctx context.Context, key string, increment int64) (int64, error) {
	r := p.c.Do(ctx, p.c.B().Incrby().Key(key).Increment(increment).Build())
	p.l.WithContext(ctx).Debug("Execute incrby command.", attrKey, key, "increment", increment, "result", r.String())
	return r.AsInt64()
}
