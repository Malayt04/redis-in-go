package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	line = strings.TrimSuffix(line, "\n")
	line = strings.TrimSuffix(line, "\r")

	return line, nil
}

func parseBulkString(reader *bufio.Reader) (string, error) {
	lengthLine, err := readLine(reader)
	if err != nil {
		return "", err
	}

	length, err := strconv.Atoi(lengthLine)
	if err != nil {
		return "", err
	}

	if length == -1 {
		return "", nil
	}

	data := make([]byte, length)

	_, err = io.ReadFull(reader, data)
	if err != nil {
		return "", err
	}

	_, err = reader.ReadByte()
	if err != nil {
		return "", err
	}

	_, err = reader.ReadByte()
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func Decode(reader *bufio.Reader) ([]string, error) {
	prefix, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	if prefix != '*' {
		return nil, fmt.Errorf("expected array")
	}

	countLine, err := readLine(reader)
	if err != nil {
		return nil, err
	}

	count, err := strconv.Atoi(countLine)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, count)

	for i := 0; i < count; i++ {
		elemType, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}

		if elemType != '$' {
			return nil, fmt.Errorf("expected bulk string")
		}

		value, err := parseBulkString(reader)
		if err != nil {
			return nil, err
		}

		result = append(result, value)
	}

	return result, nil
}

func EncodeString(s string) string {
	return "+" + s + "\r\n"
}

func EncodeError(s string) string {
	return "-" + s + "\r\n"
}

func EncodeBulkString(s string) string {
	if s == "" {
		return "$-1\r\n"
	}
	return "$" + strconv.Itoa(len(s)) + "\r\n" + s + "\r\n"
}

func EncodeInteger(n int) string {
	return ":" + strconv.Itoa(n) + "\r\n"
}

func EncodeNull() string {
	return "$-1\r\n"
}

type Response struct {
	Type  string
	Value string
}

func DecodeResponse(reader *bufio.Reader) (Response, error) {
	prefix, err := reader.ReadByte()
	if err != nil {
		return Response{}, err
	}

	switch prefix {
	case '+':
		line, _ := readLine(reader)
		return Response{Type: "simple", Value: line}, nil
	case '-':
		line, _ := readLine(reader)
		return Response{Type: "error", Value: line}, nil
	case ':':
		line, _ := readLine(reader)
		return Response{Type: "integer", Value: line}, nil
	case '$':
		return decodeBulkString(reader)
	case '*':
		return decodeArray(reader)
	default:
		return Response{}, fmt.Errorf("unknown response type: %c", prefix)
	}
}

func decodeBulkString(reader *bufio.Reader) (Response, error) {
	lengthLine, err := readLine(reader)
	if err != nil {
		return Response{}, err
	}

	length, err := strconv.Atoi(lengthLine)
	if err != nil {
		return Response{}, err
	}

	if length == -1 {
		return Response{Type: "null", Value: ""}, nil
	}

	data := make([]byte, length)
	_, err = io.ReadFull(reader, data)
	if err != nil {
		return Response{}, err
	}

	reader.ReadByte()
	reader.ReadByte()

	return Response{Type: "bulk", Value: string(data)}, nil
}

func decodeArray(reader *bufio.Reader) (Response, error) {
	countLine, err := readLine(reader)
	if err != nil {
		return Response{}, err
	}

	count, err := strconv.Atoi(countLine)
	if err != nil {
		return Response{}, err
	}

	result := make([]string, 0, count)
	for i := 0; i < count; i++ {
		elemType, err := reader.ReadByte()
		if err != nil {
			return Response{}, err
		}

		if elemType == '$' {
			bs, err := decodeBulkString(reader)
			if err != nil {
				return Response{}, err
			}
			result = append(result, bs.Value)
		}
	}

	return Response{Type: "array", Value: strings.Join(result, " ")}, nil
}

func Encode(input string) string {
	parts := strings.Split(input, " ")

	var result strings.Builder

	result.WriteString("*")
	result.WriteString(strconv.Itoa(len(parts)))
	result.WriteString("\r\n")

	for _, part := range parts {
		result.WriteString("$")
		result.WriteString(strconv.Itoa(len(part)))
		result.WriteString("\r\n")
		result.WriteString(part)
		result.WriteString("\r\n")
	}

	return result.String()
}