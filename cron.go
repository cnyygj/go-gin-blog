package main

import (
	"time"
	"log"

	"github.com/robfig/cron"
	"github.com/Songkun007/go-gin-blog/models"
)

func main() {
	log.Println("Starting...")

	c := cron.New()

	// AddFunc 会向 Cron job runner 添加一个 func ，以按给定的时间表运行
	c.AddFunc("* * * * * *", func() {
		log.Println("Run models.CleanAllTag...")
		models.CleanAllTag()
	})

	c.AddFunc("* * * * * *", func() {
		log.Println("Run models.CleanAllArticle...")
		models.CleanAllArticle()
	})

	c.Start()

	t1 := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-t1.C:
			t1.Reset(time.Second * 10)
		}
	}
}