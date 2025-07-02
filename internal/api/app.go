package api

import (
	"log/slog"
	"net/http"
	"os"

	_ "example.com/m/docs"
	"example.com/m/internal/config"
	"example.com/m/internal/middleware"
	jwtService "example.com/m/internal/service/jwt"
	userService "example.com/m/internal/service/user"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type State struct {
	AuthService *userService.AuthService
	JWTService  *jwtService.JWTService
}

func New(config config.Config, pool *pgxpool.Pool) http.Handler {
	state := State{
		AuthService: &userService.AuthService{Pool: pool},
		JWTService:  &jwtService.JWTService{Salt: config.SecretSalt, Pool: pool, WebhookURL: config.WebhookURL},
	}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Healthy"))
	})
	mux.HandleFunc("/docs/", httpSwagger.Handler(
		httpSwagger.URL("/docs/doc.json"),
		httpSwagger.DeepLinking(true),
	))
	apiMux := http.NewServeMux()

	apiMux.HandleFunc("POST /login", state.LoginHandler)
	apiMux.HandleFunc("POST /refresh", state.RefreshHandler)
	apiMux.HandleFunc("POST /logout", middleware.JWTMiddleware(state.JWTService)(state.LogoutHandler))
	apiMux.HandleFunc("GET /me", middleware.JWTMiddleware(state.JWTService)(state.MeHandler))

	mux.Handle("/api/", http.StripPrefix("/api", apiMux))

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	finalMux := middleware.LoggerMiddleware(logger, mux)
	return finalMux
}
