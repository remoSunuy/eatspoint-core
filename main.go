package main

import (
	"fmt"
	"time"

	"github.com/remoSunuy/eatspoint-core/consul"
)

func main() {
	for i := 1; i <= 100; i++ {
    	userPoint, err := consul.HealthConsulService("eatspoint-user","")
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("User Service URL " + userPoint)
		time.Sleep(5 * time.Second)
	}
}