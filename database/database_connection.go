package database

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DBPool *pgxpool.Pool

func Connection() {
	ctx := context.Background()
	err := godotenv.Load(".env")

	if err != nil {
		log.Printf("No .env file loaded, using runtime environment variables: %v", err)
	}

	DATABASE_URL := os.Getenv("DATABASE_URL")

	if DATABASE_URL == "" {
		log.Fatal("DATABASE_URL is not configured")
	}

	MAX_CONNS := os.Getenv("MAX_CONNS")

	if MAX_CONNS == "" {
		log.Fatal("MAX_CONNS is not configured")
	}

	MIN_CONNS := os.Getenv("MIN_CONNS")

	if MIN_CONNS == "" {
		log.Fatal("MIN_CONNS is not configured")
	}

	config, err := pgxpool.ParseConfig(DATABASE_URL)

	if err != nil {
		log.Fatal(err)
	}

	max_conn, err := strconv.ParseInt(MAX_CONNS, 10, 32)

	if err != nil {
		log.Fatal(err)
	}

	min_conn, err := strconv.ParseInt(MIN_CONNS, 10, 32)

	if err != nil {
		log.Fatal(err)
	}

	config.MaxConns = int32(max_conn)
	config.MinConns = int32(min_conn)
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	dbpool, err := pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		log.Fatal("Unable to create connection pool:", err)
	}

	DBPool = dbpool

	err = dbpool.Ping(ctx)

	if err != nil {
		log.Fatal("Database ping failed:", err)
	}

	log.Println("Connected to PostgreSQL successfully")

}
