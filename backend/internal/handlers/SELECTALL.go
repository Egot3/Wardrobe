package handlers

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetSelectedAll(dbpool *pgxpool.Pool, raw string) ([]map[string]interface{}, error) {
	query := string(raw)

	rows, err := dbpool.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("db query failed")
	}
	defer rows.Close()

	var results []map[string]interface{}

	for rows.Next() {
		FieldDescriptions := rows.FieldDescriptions()
		values, err := rows.Values()
		if err != nil {
			return nil, fmt.Errorf("rows' values fetching failed: %v", err.Error())
		}

		rowMap := make(map[string]interface{})
		for i, fd := range FieldDescriptions {
			rowMap[string(fd.Name)] = values[i]
		}
		results = append(results, rowMap)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row's iteration failed: %v", err.Error())
	}

	return results, nil

}
