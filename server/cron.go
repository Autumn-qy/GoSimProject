package server

import "github.com/robfig/cron/v3"

func Cron() {
	c := cron.New()

	_, _ = c.AddFunc("0/10 * * * *", FetchConfig)
	_, _ = c.AddFunc("0/10 * * * *", BatchActive)

	c.Start()
	select {}
}
