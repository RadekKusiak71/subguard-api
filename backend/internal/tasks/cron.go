package tasks

import "github.com/robfig/cron/v3"

func StartCron(subCron *SubscriptionCron) *cron.Cron {
	cr := cron.New(cron.WithSeconds())

	// Daily at 9 AM
	_, err := cr.AddFunc("0 0 9 * * *", subCron.StartReminderCron())
	if err != nil {
		panic(err)
	}

	cr.Start()
	return cr
}
