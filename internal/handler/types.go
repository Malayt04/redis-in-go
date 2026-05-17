package handler

import (
	"strconv"

	"github.com/Malayt04/redis-in-go/internal/store"
)

type CommandHandler func(args []string, s *store.Store) Response

type ResponseType int

const (
	ResponseSimpleString ResponseType = iota
	ResponseError
	ResponseBulkString
	ResponseInteger
	ResponseNull
)

type Response struct {
	Type  ResponseType
	Value string
}

func ErrorResponse(msg string) Response {
	return Response{Type: ResponseError, Value: msg}
}

func SimpleStringResponse(msg string) Response {
	return Response{Type: ResponseSimpleString, Value: msg}
}

func BulkStringResponse(msg string) Response {
	return Response{Type: ResponseBulkString, Value: msg}
}

func IntegerResponse(n int) Response {
	return Response{Type: ResponseInteger, Value: strconv.Itoa(n)}
}

func NullResponse() Response {
	return Response{Type: ResponseNull, Value: ""}
}