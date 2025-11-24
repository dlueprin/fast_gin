package cron_ser

import (
	"fmt"
	"time"
)

func HelloCron() {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("I am waiting for you!")
}
