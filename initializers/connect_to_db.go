package initializers

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectToDb() {
	if DB != nil {
		DB.Close()
	}

	connStr := os.Getenv("DB")

	conn, err := pgxpool.New(context.Background(), connStr)

	if err != nil {
		log.Fatalf("Erreur de connexion DB: %v\n", err)
	}
	fmt.Println("✅ Connexion DB réussie !")

	if err != nil {
		panic(err)
	}
	DB = conn
}
