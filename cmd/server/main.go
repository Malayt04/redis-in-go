package main

import (
	"github.com/Malayt04/redis-in-go/internal/server"
)

func main() {
	s := server.New(":6379")
	s.Start()
}