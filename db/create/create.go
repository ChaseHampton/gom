package create

import (
	"context"
	"log"
	"time"

	"github.com/chasehampton/gom/db"
	"github.com/chasehampton/gom/models"
)

type LogInserter interface {
	InsertLogEntry(ctx context.Context, entry models.LogEntry) error
}

type InserterImpl struct {
}

func InsertProject(ctx context.Context, project models.Project) error {
	query := "call usp_insert_project($1, $2, $3, $4)"
	_, err := db.DBPool.Exec(ctx, query, project.Name, project.Description, time.Now(), db.GetCurrentUser())
	if err != nil {
		log.Printf("Error inserting project: %v\n", err)
		return err
	}
	return nil
}

func InsertConnection(ctx context.Context, connection models.Connection) error {
	query := "call usp_insert_connection($1, $2, $3, $4)"
	_, err := db.DBPool.Exec(ctx, query, connection.Protocol.ProtocolID, connection.AuthDetail.AuthID, time.Now(), db.GetCurrentUser())
	if err != nil {
		log.Printf("Error inserting connection: %v\n", err)
		return err
	}
	return nil
}

func InsertProtocol(ctx context.Context, protocol models.Protocol) error {
	query := "call usp_insert_protocol($1, $2, $3, $4)"
	_, err := db.DBPool.Exec(ctx, query, protocol.Name, protocol.Description, time.Now(), db.GetCurrentUser())
	if err != nil {
		log.Printf("Error inserting protocol: %v\n", err)
		return err
	}
	return nil
}

func InsertAction(ctx context.Context, action models.Action) error {
	query := "call usp_insert_action($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err := db.DBPool.Exec(ctx, query, action.LocalPath, action.RemotePath, action.Bucket, action.IsUpload, action.Connection.ConnectionID, action.Project.ProjectID, time.Now(), db.GetCurrentUser())
	if err != nil {
		log.Printf("Error inserting action: %v\n", err)
		return err
	}
	return nil
}

func InsertConnectionConfig(ctx context.Context, config models.ConnectionConfig, connId int) error {
	query := "call usp_insert_connection_config($1, $2, $3, $4, $5)"
	_, err := db.DBPool.Exec(ctx, query, config.Key, config.Value, connId, time.Now(), db.GetCurrentUser())
	if err != nil {
		log.Printf("Error inserting connection config: %v\n", err)
		return err
	}
	return nil
}

func InsertAuthData(ctx context.Context, auth models.AuthDetail) error {
	query := "call usp_insert_auth_detail($1, $2, $3, $4)"
	_, err := db.DBPool.Exec(ctx, query, auth.Description, auth.VaultPath, time.Now(), db.GetCurrentUser())
	if err != nil {
		log.Printf("Error inserting auth data: %v\n", err)
		return err
	}
	return nil
}

func InsertSchedule(ctx context.Context, schedule models.Schedule) error {
	query := "call usp_insert_schedule($1, $2, $3)"
	_, err := db.DBPool.Exec(ctx, query, schedule.Description, time.Now(), db.GetCurrentUser())
	if err != nil {
		log.Printf("Error inserting schedule: %v\n", err)
		return err
	}
	return nil
}

func InsertScheduleAction(ctx context.Context, scheduleAction models.ScheduleAction) error {
	query := "call usp_insert_schedule_action($1, $2)"
	_, err := db.DBPool.Exec(ctx, query, scheduleAction.ScheduleID, scheduleAction.ActionID)
	if err != nil {
		log.Printf("Error inserting schedule action: %v\n", err)
		return err
	}
	return nil
}

func InsertLogEntry(ctx context.Context, logEntry models.LogEntry) error {
	query := "call usp_insert_log_entry($1, $2, $3, $4, $5)"
	_, err := db.DBPool.Exec(ctx, query, logEntry.Message, logEntry.StackTrace, logEntry.ActionID, logEntry.Additional)
	if err != nil {
		log.Printf("Error inserting log entry: %v\n", err)
		return err
	}
	return nil
}

func (i *InserterImpl) InsertLogEntry(ctx context.Context, entry models.LogEntry) error {
	return InsertLogEntry(ctx, entry)
}
