package request

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"unicode"
)

type ParserState int

const (
	StateInit ParserState = iota
	StateDone
)

type Request struct {
	RequestLine RequestLine
	parserState ParserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

// It accepts the next slice of bytes that needs to be parsed into the Request struct
// It updates the "state" of the parser, and the parsed RequestLine field.
// It returns the number of bytes it consumed (meaning successfully parsed) and an error if it encountered one.
func (r *Request) parse(data []byte) (int, error) {
	switch r.parserState {
	case StateInit:
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		//n == 0 means parseRequestLine needs more data to find /r/n
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.parserState = StateDone
		return n, nil
	case StateDone:
		return 0, errors.New("parser in state Done. No more read allowed")
	default:
		return 0, errors.New("unknown parser state")
	}
}

// GET /coffee HTTP/1.1
// Host: localhost:42069
// User-Agent: curl/7.81.0
// Accept: */*
//
// {"flavor":"dark mode"}
func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, 8)
	readToIndex := 0
	req := &Request{
		parserState: StateInit,
	}
	for req.parserState != StateDone {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		numBytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				req.parserState = StateDone
				break
			}
			return nil, err
		}
		readToIndex += numBytesRead

		numBytesParsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[numBytesParsed:])
		readToIndex -= numBytesParsed
	}
	return req, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	i := bytes.Index(data, []byte("\r\n"))
	if i == -1 {
		return nil, 0, nil
	}
	requestLineText := string(data[:i])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}
	return requestLine, i, nil
}

func requestLineFromString(data string) (*RequestLine, error) {
	parts := strings.Split(data, " ") //should now have method at 0, requesttarget at 1 and version at 2
	if len(parts) != 3 {
		return nil, errors.New("incorrect number of fields in RequestLine")
	}
	if parts[2] != "HTTP/1.1" {
		return nil, errors.New("unsupported or incorrect HTTP Version")
	}
	version := strings.Split(parts[2], "/")

	for _, r := range parts[0] {
		if !unicode.IsUpper(r) {
			return nil, errors.New("unsupported method")
		}
	}

	return &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   version[1],
	}, nil
}
