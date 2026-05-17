package handler

import (
	"strconv"
	"time"

	"github.com/Malayt04/redis-in-go/internal/store"
)

func setCommand(args []string, s *store.Store) Response {
	if len(args) < 2 {
		return ErrorResponse("ERR wrong number of arguments for 'SET' command")
	}

	key := args[0]
	value := args[1]
	expiry := time.Duration(0)

	// Check for EX option
	for i := 2; i < len(args)-1; i++ {
		if args[i] == "ex" {
			sec, err := strconv.Atoi(args[i+1])
			if err != nil {
				return ErrorResponse("ERR invalid expiry time")
			}
			expiry = time.Duration(sec) * time.Second
			break
		}
	}

	if expiry > 0 {
		s.SetWithExpiry(key, value, expiry)
	} else {
		s.Set(key, value)
	}

	return SimpleStringResponse("OK")
}