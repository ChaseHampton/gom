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
	currentUser string
	once        sync.Once
)

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
	once.Do(func() {
		cUser, err := user.Current()
		if err != nil {
			currentUser = "unknown"
		} else {
			currentUser = cUser.Username
		}
	})
	return currentUser
}
