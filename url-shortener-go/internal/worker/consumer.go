package worker

import (
	"database/sql"
	"fmt"
	"url-shortener/internal/stats"
)

func StartConsumer(db *sql.DB, events chan string) {
	go func() {
		fmt.Println("Stats consumer started")
		for code := range events {
			err := stats.InsertClick(db, code)
			if err != nil {
				fmt.Println("Failed stats insert:", err)
			}
		}

	}()
}
