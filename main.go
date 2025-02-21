package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/andrefelizardo/posgoexpert_rate-limiter/limiter"
	"github.com/andrefelizardo/posgoexpert_rate-limiter/middleware"
	"github.com/andrefelizardo/posgoexpert_rate-limiter/persistence"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Nenhum .env encontrado, usando variáveis de ambiente")
	}

	ipLimit, _ := strconv.Atoi(os.Getenv("RATE_LIMIT_IP"))
	tokenLimit, _ := strconv.Atoi(os.Getenv("RATE_LIMIT_TOKEN"))
	blockTimeIP, _ := strconv.Atoi(os.Getenv("BLOCK_TIME_IP"))
	blockTimeToken, _ := strconv.Atoi(os.Getenv("BLOCK_TIME_TOKEN"))

	store, err := persistence.NewRedisStore(os.Getenv("REDIS_ADDR"))
	if err != nil {
		log.Fatalf("Falha ao conectar ao Redis: %v", err)
	}

	rl := limiter.NewRateLimiter(ipLimit, tokenLimit, blockTimeIP, blockTimeToken, store)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	handler := middleware.RateLimiterMiddleware(rl)(mux)

	log.Println("Servidor executando na porta :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Erro no servidor: %v", err)
	}
}
