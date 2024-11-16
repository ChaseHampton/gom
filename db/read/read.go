package read

import (
	"context"

	"github.com/chasehampton/gom/db"
	"github.com/chasehampton/gom/models"
)

func GetSchedule(ctx context.Context, scheduleID int) ([]models.Schedule, error) {
	query := "select * FROM usp_collect_schedule_actions($1)"
	rows, err := db.DBPool.Query(ctx, query, scheduleID)
	if err != nil {
		return []models.Schedule{}, err
	}
	defer rows.Close()

	schedulesMap := make(map[int]*models.Schedule, 0)
	for rows.Next() {
		var (
			s  models.Schedule
			a  models.Action
			c  models.Connection
			p  models.Project
			pr models.Protocol
			cc models.ConnectionConfig
			ad models.AuthDetail
		)

		err := rows.Scan(
			&s.ScheduleID,
			&s.Description,
			&s.DateAdded,
			&s.AddedBy,
			&p.ProjectID,
			&p.Name,
			&p.Description,
			&p.DateAdded,
			&p.AddedBy,
			&a.ActionID,
			&a.LocalPath,
			&a.RemotePath,
			&a.Bucket,
			&a.DateAdded,
			&a.AddedBy,
			&c.ConnectionID,
			&c.DateAdded,
			&c.AddedBy,
			&pr.ProtocolID,
			&pr.Name,
			&pr.Description,
			&pr.DateAdded,
			&pr.AddedBy,
			&cc.ConfigID,
			&cc.Key,
			&cc.Value,
			&cc.DateAdded,
			&cc.AddedBy,
			&ad.AuthID,
			&ad.Description,
			&ad.Username,
			&ad.DateAdded,
			&ad.AddedBy,
		)

		if err != nil {
			return []models.Schedule{}, err
		}

		if existSched, found := schedulesMap[s.ScheduleID]; found {
			a.Project = p
			a.Connection = c
			c.Protocol = pr
			c.Configs = append(c.Configs, cc)
			c.AuthDetail = ad
			existSched.Actions = append(existSched.Actions, a)
		} else {
			a.Project = p
			a.Connection = c
			c.Protocol = pr
			if cc.ConfigID.Valid {
				c.Configs = []models.ConnectionConfig{cc}
			} else {
				c.Configs = []models.ConnectionConfig{} // Create an empty slice
			}
			c.AuthDetail = ad
			s.Actions = []models.Action{a}
			schedulesMap[s.ScheduleID] = &s
		}
	}
	result := make([]models.Schedule, 0)
	for _, sched := range schedulesMap {
		result = append(result, *sched)
	}
	return result, nil
}
