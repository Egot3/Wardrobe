package main

import (
	"context"
	"fmt"
	"os"
	"wardrobe/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	//env

	if err := godotenv.Load(); err != nil {
		panic("no .env")
	}

	connectionString, exists := os.LookupEnv("DATABASE_URL")
	if !exists {
		panic("no db")
	}

	port, exists := os.LookupEnv("SERVER_PORT")
	if !exists {
		fmt.Print("no port var found, switching to 8080")
		port = "8080"
	}

	//sql part

	dbpool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		panic(fmt.Sprintf("unableToConnect: %v\n", err))
	}
	defer dbpool.Close()

	err = dbpool.Ping(context.Background())
	if err != nil {
		panic(fmt.Sprintf("cannotPingDB %v\n", err))
	}
	fmt.Println("connectedToPSQL!")

	//gin's part
	router := gin.Default()
	router.GET("/query", handlers.Dispatcher(dbpool))
	router.Run("localhost:" + port)
}
