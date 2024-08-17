package market

import (
	"omni-balance/internal/daemons"
	"time"
)

func init() {
	daemons.RegisterIntervalTask(daemons.Task{
		Name:            "market",
		Description:     "Look for tasks in the database that are not being processed and process them.",
		TaskFunc:        Run,
		DefaultInterval: time.Minute * 3,
		RunOnStart:      false,
	})
}
