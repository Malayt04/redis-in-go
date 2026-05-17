package handler

import (
	"strings"

	"github.com/Malayt04/redis-in-go/internal/store"
)

var commands = map[string]CommandHandler{
	"ping":  pingCommand,
	"set":   setCommand,
	"get":   getCommand,
	"del":   delCommand,
	"exists": existsCommand,
	"expire": expireCommand,
	"ttl":   ttlCommand,
}

type Handler struct {
	store *store.Store
}

func New(s *store.Store) *Handler {
	return &Handler{store: s}
}

func (h *Handler) Handle(message []string) Response {
	if len(message) == 0 {
		return ErrorResponse("ERR empty command")
	}

	command := strings.ToLower(message[0])
	args := message[1:]

	handler, ok := commands[command]
	if !ok {
		return ErrorResponse("ERR unknown command '" + command + "'")
	}

	return handler(args, h.store)
}