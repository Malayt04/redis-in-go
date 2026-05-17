package handler

import "github.com/Malayt04/redis-in-go/internal/store"

func delCommand(args []string, s *store.Store) Response {
	if len(args) != 1 {
		return ErrorResponse("ERR wrong number of arguments for 'DEL' command")
	}

	s.Delete(args[0])
	return SimpleStringResponse("OK")
}