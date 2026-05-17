package handler

import "github.com/Malayt04/redis-in-go/internal/store"

func getCommand(args []string, s *store.Store) Response {
	if len(args) != 1 {
		return ErrorResponse("ERR wrong number of arguments for 'GET' command")
	}

	val, ok := s.Get(args[0])
	if !ok {
		return NullResponse()
	}

	return BulkStringResponse(val)
}