package db

import (
	"context"
	"fmt"
	"log"
	"os/user"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	DBPool      *pgxpool.Pool
	CurrentUser string
	O           sync.Once
)

var CurrentFunc = user.Current

func Connect(connStr string) {
	var err error
	DBPool, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	fmt.Println("Connected to database")
}

func Close() {
	DBPool.Close()
}

func GetCurrentUser() string {
	O.Do(func() {
		cUser, err := CurrentFunc()
		if err != nil {
			CurrentUser = "unknown"
		} else {
			CurrentUser = cUser.Username
		}
	})
	return CurrentUser
}
