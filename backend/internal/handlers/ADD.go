package handlers

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AddData(dbpool *pgxpool.Pool, parced []string) error {
	if len(parced) < 3 {
		return fmt.Errorf("not enough arguments")
	}
	whatTo := parced[1] //unintentional reference

	transaction, err := dbpool.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("could not begin a transaction: %v", err.Error())
	}
	defer transaction.Rollback(context.Background())

	if whatTo == "file" {
		path := parced[2]
		fileInfo, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("file doesn't seem to exist: %v", err.Error())
		}

		fileName := fileInfo.Name()

		var fileId string
		err = transaction.QueryRow(context.Background(),
			`INSERT INTO files (name, storage_path) VALUES ($1, $2) RETURNING id`, fileName, path,
		).Scan(&fileId) //link-like goes hard
		if err != nil {
			return fmt.Errorf(`file is already stored: %v`, err.Error())
		}

		if len(parced) >= 4 { //living hell of an alternative query
			pipeCommand := parced[3]
			switch pipeCommand {
			case "WITH":
				argument := parced[4]
				switch argument {
				case "tag":

					tagName := parced[5]

					var tagId string
					err := transaction.QueryRow(context.Background(), fmt.Sprintf(`INSERT INTO tags (name) Values ('%v') ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id`, tagName)).Scan(&tagId)
					if err != nil {
						return fmt.Errorf("couldn't check if tag exists(or create one)")
					}

					_, err = transaction.Exec(context.Background(),
						`INSERT INTO file_tags VALUES ($1,$2)`, fileId, tagId,
					)

					if err != nil {
						return fmt.Errorf("additional query failed: %v", err.Error())
					}

				}

			}
		}
		if err := transaction.Commit(context.Background()); err != nil {
			return fmt.Errorf("transaction wasn't commited: %v", err.Error())
		}
	}

	return nil
}
