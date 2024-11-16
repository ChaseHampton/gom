package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/chasehampton/gom/db"
	"github.com/chasehampton/gom/db/read"
	"github.com/chasehampton/gom/handlers"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)
	}
	connStr := os.Getenv("GOM_DB_CONN_STR")

	schedule_id := flag.Int("schedule_id", 1, "The schedule ID being requested")
	flag.Parse()

	if *schedule_id == 0 && flag.NArg() > 0 {
		arg := flag.Arg(0)
		id, err := strconv.Atoi(arg)
		if err != nil {
			log.Fatalf("Invalid schedule ID passed to args: %v\n", arg)
		}
		*schedule_id = id
	}

	db.Connect(connStr)
	defer db.Close()

	schedule, err := read.GetSchedule(context.Background(), *schedule_id)
	if err != nil {
		log.Fatalf("Error getting schedule: %v\n", err)
	}

	bhandler := handlers.GetBaseHandler()
	for _, s := range schedule {
		for _, a := range s.Actions {
			if !a.LocalPath.Valid {
				log.Printf("Local path is not valid for action: %v\n", a.ActionID)
			}
			files, err := bhandler.GetUploadFiles(a.LocalPath.String)
			if err != nil {
				log.Fatalf("Error getting files: %v\n", err)
			}
			for _, f := range files {
				fmt.Println(f.Name())
			}
		}
	}
}
