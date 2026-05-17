package handler

import "github.com/Malayt04/redis-in-go/internal/store"

func ttlCommand(args []string, s *store.Store) Response {
	if len(args) != 1 {
		return ErrorResponse("ERR wrong number of arguments for 'TTL' command")
	}

	ttl := s.TTL(args[0])
	return IntegerResponse(ttl)
}