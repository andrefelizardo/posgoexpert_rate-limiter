package integration

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/andrefelizardo/posgoexpert_rate-limiter/limiter"
	"github.com/andrefelizardo/posgoexpert_rate-limiter/middleware"
	"github.com/andrefelizardo/posgoexpert_rate-limiter/persistence"
)

func TestIntegrationRateLimiter(t *testing.T) {
	// Carrega configurações via ambiente (ou use valores fixos para os testes)
	ipLimit := 3
	tokenLimit := 5
	blockTimeIP := 60    // segs
	blockTimeToken := 60 // segs

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	store, err := persistence.NewRedisStore(redisAddr)
	if err != nil {
		t.Fatalf("Falha ao conectar ao Redis: %v", err)
	}

	rl := limiter.NewRateLimiter(ipLimit, tokenLimit, blockTimeIP, blockTimeToken, store)
	mux := http.NewServeMux()
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	handler := middleware.RateLimiterMiddleware(rl)(mux)
	server := httptest.NewServer(handler)
	defer server.Close()

	client := server.Client()

	// Testa limite para IP (sem token)
	req, _ := http.NewRequest("GET", server.URL+"/test", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	var lastResp *http.Response
	for i := 1; i <= ipLimit+1; i++ {
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Erro na requisição %d: %v", i, err)
		}
		lastResp = resp
	}
	body, _ := ioutil.ReadAll(lastResp.Body)
	lastResp.Body.Close()
	if lastResp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Requisição excedente para IP: esperado 429, obtido %d com resposta: %s", lastResp.StatusCode, string(body))
	}

	// Testa limite para Token (o token se sobrepõe ao IP)
	req, _ = http.NewRequest("GET", server.URL+"/test", nil)
	req.Header.Set("API_KEY", "abc123")
	req.RemoteAddr = "127.0.0.1:1234"
	for i := 1; i <= tokenLimit+1; i++ {
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Erro na requisição token %d: %v", i, err)
		}
		lastResp = resp
	}
	body, _ = ioutil.ReadAll(lastResp.Body)
	lastResp.Body.Close()
	if lastResp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Requisição excedente para token: esperado 429, obtido %d com resposta: %s", lastResp.StatusCode, string(body))
	}
}
