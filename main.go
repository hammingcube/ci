package main

import (
	"fmt"
	"github.com/phayes/hookserve/hookserve"
	"time"
)

func main() {
	server := hookserve.NewServer()
	server.Port = 8120
	server.Secret = "absolutesecret"
	server.GoListenAndServe()
	fmt.Printf("Listening on %d\n", server.Port)
	for {
		select {
		case commit := <-server.Events:
				fmt.Println(commit.Owner, commit.Repo, commit.Branch, commit.Commit)		
		default:
			time.Sleep(100)
			//fmt.Println("No activity...")
		}
	}
}
