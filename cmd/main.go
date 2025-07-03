// @title						AuthApi
// @version					1.0
// @description				For trials
// @host						localhost:8000
// @BasePath					/
// @securityDefinitions.apikey	Bearer
// @authorizationurl			http://localhost:8000/api/login
// @in							header
// @name						Authorization
// @description				Paste token with Bearer prefix, e.g. `Bearer <your-token>`
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fedya-eremin/medods-trials/internal/api"
	"github.com/fedya-eremin/medods-trials/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config := config.Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		SecretSalt:  os.Getenv("SECRET_SALT"),
		WebhookURL:  os.Getenv("WEBHOOK_URL"),
	}
	dbconfig, err := pgxpool.ParseConfig(config.DatabaseURL)
	if err != nil {
		log.Fatalln("DATABASE_URL is wrong")
		os.Exit(1)
	}
	pool, err := pgxpool.NewWithConfig(
		context.Background(),
		dbconfig,
	)
	if err != nil {
		fmt.Println("Failed to connect to DB")
		os.Exit(2)
	}
	mux := api.New(config, pool)
	fmt.Println("=== Server Started at :8000 ===")
	fmt.Println("=== Swagger is at /docs/index.html ===")
	if err := http.ListenAndServe(":8000", mux); err != nil {
		log.Fatalln("Something went wrong")
		os.Exit(3)
	}
}
