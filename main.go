package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/chasehampton/gom/db"
	"github.com/chasehampton/gom/db/create"
	"github.com/chasehampton/gom/db/read"
	"github.com/chasehampton/gom/handlers"
	"github.com/chasehampton/gom/logger"
	"github.com/chasehampton/gom/models"
	"github.com/chasehampton/gom/runner"
	"github.com/chasehampton/gom/vaultclient"
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

	l := logger.NewLogger(&create.InserterImpl{})

	schedule, err := read.GetSchedule(context.Background(), *schedule_id)
	if err != nil {
		log.Fatalf("Error getting schedule: %v\n", err)
	}

	var wg sync.WaitGroup

	for _, s := range schedule {
		if s.Actions == nil {
			log.Printf("No actions found for schedule: %v\n", s.ScheduleID)
			continue
		}
		vc, err := vaultclient.NewVaultClient(os.Getenv("VAULT_ADDR"), os.Getenv("VAULT_DEV_ROOT_TOKEN_ID"))
		if err != nil {
			log.Printf("Error creating vault client: %v\n", err)
		}
		bh := handlers.GetBaseHandler(vc, l)
		for _, a := range s.Actions {
			wg.Add(1)
			go func(act models.Action) {
				defer wg.Done()
				log.Printf("Action: %v\n", act)
				err := runner.Run(act, bh)
				if err != nil {
					log.Printf("Error running action: %v\n", err)
				}
			}(a)
		}
	}
	wg.Wait()
}
