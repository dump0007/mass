package services

import (
	"fmt"
	"mass/helpers"
	"time"

	"github.com/robfig/cron/v3"
)

func Cron(){
	c := cron.New(cron.WithSeconds())
	
	// Add a job to run every minute
	c.AddFunc("*/30 * * * * *", func() {
		fmt.Println("Job running every minute:", time.Now())
		helpers.ExecuteOrderProcessing()
	})

	// // Add a job to run every day at midnight
	// c.AddFunc("0 0 0 * * *", func() {
	// 	fmt.Println("Job running every day at midnight:", time.Now())
	// })

	// // Add a job to run every Monday at 9am
	// c.AddFunc("0 0 9 * * 1", func() {
	// 	fmt.Println("Job running every Monday at 9am:", time.Now())
	// })

	// Start the cron scheduler
	c.Start()

	// Keep the program running
	select {}

}