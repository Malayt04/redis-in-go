package handler

import "github.com/Malayt04/redis-in-go/internal/store"

func pingCommand(args []string, s *store.Store) Response {
	return SimpleStringResponse("PONG")
}