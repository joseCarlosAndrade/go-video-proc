package main

import (
	// "fmt"

	"database/sql" // only the sql interface, does not implement all dialects
	"fmt"
	"log/slog"
	"os"
	"videoproc/internal/converter"
	_ "github.com/lib/pq" // postgres implementation to use alongside database/sql
)

func connectPostgres() (*sql.DB, error) {
	user := getEnvOrDefault("POSTGRES_USER", "user")
	password := getEnvOrDefault("POSTGRES_PASSWORD", "password")
	dbname := getEnvOrDefault("POSTGRES_DB", "converter")
	host := getEnvOrDefault("POSTGRES_HOST", "postgres")
	sslmode := getEnvOrDefault("POSTGRES_SSLMODE", "disable")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=%s", user, password, dbname, host, sslmode)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		slog.Error("error connecting to database", slog.String("connStr", connStr))
		return nil, err
	}
	
	err = db.Ping()
	if err != nil {
		slog.Error("error pinging database", slog.String("connStr", connStr))
		return nil, err
	}
	
	slog.Info("successfully connected to postgres")

	return db, nil
}

// tries to read an env value. if doesnt exist, use default value
func getEnvOrDefault(key, defaultValue string) string {
	if value, exist := os.LookupEnv(key); exist {
		return value
	}
	return defaultValue
}

func main() {
	// mergeChunks("mediatest/media/uploads/1", "merged.mp4" )
	// vc := NewVideo
	db, err := connectPostgres()

	if err != nil {
		panic(err)
	}

	vc := converter.NewVideoConverter(db)
	vc.Handle([]byte(`{ "video_id" : 1 , "path" : "/media/uploads/1" }`))
}

