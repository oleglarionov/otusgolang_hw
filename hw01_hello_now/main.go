package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	currentTime := time.Now()
	exactTime, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("current time: %s\n", currentTime.Round(0))
	fmt.Printf("exact time: %s\n", exactTime.Round(0))
}
