package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Dispatcher(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, _ := c.GetRawData()
		userInput := strings.TrimRight(strings.TrimSpace(string(raw)), ";")

		if userInput == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No query"})
			return
		}

		parced := strings.Fields(userInput)
		command := strings.ToUpper(parced[0])

		switch command {
		case "SELECT":
			data, err := GetSelectedAll(dbpool, userInput)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
		case "ADD":
			err := AddData(dbpool, parced)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.Status(http.StatusNoContent)
		}
	}
}
