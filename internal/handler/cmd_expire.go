package handler

import (
	"strconv"
	"time"

	"github.com/Malayt04/redis-in-go/internal/store"
)

func expireCommand(args []string, s *store.Store) Response {
	if len(args) != 2 {
		return ErrorResponse("ERR wrong number of arguments for 'EXPIRE' command")
	}

	sec, err := strconv.Atoi(args[1])
	if err != nil {
		return ErrorResponse("ERR invalid expiry time")
	}

	if s.Expire(args[0], time.Duration(sec)*time.Second) {
		return IntegerResponse(1)
	}
	return IntegerResponse(0)
}