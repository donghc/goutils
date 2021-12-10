package runner

import (
	"log"
	"os"
	"testing"
	"time"
)

const timeout = 3 * time.Second

func TestRunner_Start(t *testing.T) {
	log.Println("starting work . ")

	r := New(timeout)
	r.Add(createTask(), createTask(), createTask())

	if err := r.Start(); err != nil {
		switch err {
		case ErrTimeOut:
			log.Println("Terminating due to timeout.")
			os.Exit(1)
		case ErrInterrupt:
			log.Println("Terminating due to interrupt.")
		}
	}

	log.Println("Process ended.")
}

func createTask() func(int) {
	return func(id int) {
		log.Printf("Processor -> Task #%d.", id)
		time.Sleep(time.Duration(id) * time.Second)
	}
}
