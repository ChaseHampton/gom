package models

import (
	"database/sql"
	"time"
)

// Project represents a project record in the database
type Project struct {
	ProjectID   int            `db:"project_id" json:"project_id"`
	Name        sql.NullString `db:"name" json:"name"`
	Description sql.NullString `db:"description" json:"description"`
	DateAdded   time.Time      `db:"date_added" json:"dateAdded"`
	AddedBy     sql.NullString `db:"added_by" json:"addedBy"`
}

// Connection represents a connection record in the database
type Connection struct {
	ConnectionID int                `db:"connection_id" json:"connection_id"`
	Protocol     Protocol           `json:"protocol"`
	AuthDetail   AuthDetail         `json:"auth_detail"`
	DateAdded    time.Time          `db:"date_added" json:"dateAdded"`
	AddedBy      sql.NullString     `db:"added_by" json:"addedBy"`
	Configs      []ConnectionConfig `json:"configs"`
}

// Protocol represents a protocol record in the database
type Protocol struct {
	ProtocolID  int            `db:"protocol_id" json:"protocol_id"`
	Name        sql.NullString `db:"name" json:"name"`
	Description sql.NullString `db:"description" json:"description"`
	DateAdded   time.Time      `db:"date_added" json:"dateAdded"`
	AddedBy     sql.NullString `db:"added_by" json:"addedBy"`
}

// Action represents an action record in the database
type Action struct {
	ActionID   int            `db:"action_id" json:"action_id"`
	LocalPath  sql.NullString `db:"local_path" json:"localPath"`
	RemotePath sql.NullString `db:"remote_path" json:"remotePath"`
	Bucket     sql.NullString `db:"bucket" json:"bucket"`
	Connection Connection     `db:"connection_id" json:"connection_id"`
	Project    Project        `db:"project_id" json:"project_id"`
	DateAdded  time.Time      `db:"date_added" json:"dateAdded"`
	AddedBy    sql.NullString `db:"added_by" json:"addedBy"`
}

// ConnectionConfig represents a connection configuration record in the database
type ConnectionConfig struct {
	ConfigID  sql.NullInt64  `db:"config_id" json:"config_id"`
	Key       sql.NullString `db:"key" json:"key"`
	Value     sql.NullString `db:"value" json:"value"`
	DateAdded sql.NullTime   `db:"date_added" json:"dateAdded"`
	AddedBy   sql.NullString `db:"added_by" json:"addedBy"`
}

// AuthDetail represents authentication details in the database
type AuthDetail struct {
	AuthID      int            `db:"auth_id" json:"auth_id"`
	Description sql.NullString `db:"description" json:"description"`
	Username    sql.NullString `db:"username" json:"username"`
	Password    sql.NullString `db:"password" json:"password"`
	PrivateKey  sql.NullString `db:"private_key" json:"privateKey"`
	AccessKey   sql.NullString `db:"access_key" json:"accessKey"`
	SecretKey   sql.NullString `db:"secret_key" json:"secretKey"`
	DateAdded   time.Time      `db:"date_added" json:"dateAdded"`
	AddedBy     sql.NullString `db:"added_by" json:"addedBy"`
}

// Schedule represents a schedule record in the database
type Schedule struct {
	ScheduleID     int            `db:"schedule_id" json:"schedule_id"`
	ScheduleNumber int            `db:"schedule_number" json:"scheduleNumber"`
	Description    sql.NullString `db:"description" json:"description"`
	DateAdded      time.Time      `db:"date_added" json:"dateAdded"`
	AddedBy        sql.NullString `db:"added_by" json:"addedBy"`
	Actions        []Action       `json:"actions"`
}

// ScheduleAction represents a link between schedules and actions in the database
type ScheduleAction struct {
	ScheduleID int `db:"schedule_id" json:"schedule_id"`
	ActionID   int `db:"action_id" json:"action_id"`
}
