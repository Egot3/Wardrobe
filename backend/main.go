package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func executeQuery(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Bad query"})
			return
		}
		query := string(raw)

		rows, err := dbpool.Query(context.Background(), query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db query failed: " + err.Error()})
			return
		}
		defer rows.Close()

		var results []map[string]interface{}

		for rows.Next() {
			FieldDescriptions := rows.FieldDescriptions()
			values, err := rows.Values()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "rows' values fetching failed: " + err.Error()})
				return
			}

			rowMap := make(map[string]interface{})
			for i, fd := range FieldDescriptions {
				rowMap[string(fd.Name)] = values[i]
			}
			results = append(results, rowMap)
		}

		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "row's iteration failed: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    results,
		})
	}
}

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
	router.GET("/query", executeQuery(dbpool))
	router.Run("localhost:" + port)
}
