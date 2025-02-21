package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/andrefelizardo/posgoexpert_rate-limiter/persistence"
)

type RateLimiter struct {
	ipLimit        int
	tokenLimit     int
	blockTimeIP    time.Duration
	blockTimeToken time.Duration
	store          persistence.Store
}

func NewRateLimiter(ipLimit, tokenLimit, blockTimeIP, blockTimeToken int, store persistence.Store) *RateLimiter {
	return &RateLimiter{
		ipLimit:        ipLimit,
		tokenLimit:     tokenLimit,
		blockTimeIP:    time.Duration(blockTimeIP) * time.Second,
		blockTimeToken: time.Duration(blockTimeToken) * time.Second,
		store:          store,
	}
}

func (r *RateLimiter) AllowRequest(ip, token string) (bool, error) {
	var key string
	var limit int
	var block time.Duration

	if token != "" {
		key = fmt.Sprintf("ratelimit:token:%s", token)
		limit = r.tokenLimit
		block = r.blockTimeToken
	} else {
		key = fmt.Sprintf("ratelimit:ip:%s", ip)
		limit = r.ipLimit
		block = r.blockTimeIP
	}

	ctx := context.Background()
	count, err := r.store.Incr(ctx, key, block)
	if err != nil {
		return false, err
	}
	if count > limit {
		return false, nil
	}
	return true, nil
}
