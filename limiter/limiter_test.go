package limiter

import (
	"context"
	"sync"
	"testing"
	"time"
)

type MockStore struct {
	counts map[string]int
	mu     sync.Mutex
}

func NewMockStore() *MockStore {
	return &MockStore{counts: make(map[string]int)}
}

func (m *MockStore) Incr(ctx context.Context, key string, expiration time.Duration) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counts[key]++
	return m.counts[key], nil
}

func TestAllowRequest(t *testing.T) {
	store := NewMockStore()
	// Limite de 2 para IP e 3 para Token, com bloqueio de 10 segundos.
	rl := NewRateLimiter(2, 3, 10, 10, store)

	// Teste para limitação por IP
	ip := "127.0.0.1"
	token := ""
	for i := 0; i < 2; i++ {
		allowed, err := rl.AllowRequest(ip, token)
		if err != nil || !allowed {
			t.Errorf("Tentativa %d: esperado permitido, obtido %v (erro: %v)", i+1, allowed, err)
		}
	}
	allowed, _ := rl.AllowRequest(ip, token)
	if allowed {
		t.Error("Esperado bloqueio para IP após limites excedidos")
	}

	// Teste para limitação por Token
	store = NewMockStore()
	rl = NewRateLimiter(2, 3, 10, 10, store)
	token = "abc123"
	for i := 0; i < 3; i++ {
		allowed, err := rl.AllowRequest(ip, token)
		if err != nil || !allowed {
			t.Errorf("Tentativa %d: esperado permitido para token, obtido %v (erro: %v)", i+1, allowed, err)
		}
	}
	allowed, _ = rl.AllowRequest(ip, token)
	if allowed {
		t.Error("Esperado bloqueio para token após limites excedidos")
	}
}
