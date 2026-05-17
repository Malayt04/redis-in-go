package handler

import "github.com/Malayt04/redis-in-go/internal/store"

func existsCommand(args []string, s *store.Store) Response {
	if len(args) != 1 {
		return ErrorResponse("ERR wrong number of arguments for 'EXISTS' command")
	}

	if s.Exists(args[0]) {
		return IntegerResponse(1)
	}
	return IntegerResponse(0)
}