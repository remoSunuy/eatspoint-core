package service

import (
	"log"
	"os"
)

func Hostname() (string) {
	hn, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}
	return hn
}