package headers

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strings"
	"unicode"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	if !haveCompleteData(data) {
		return 0, false, nil
	}
	if !verifyData(data) {
		return 0, false, errors.New("malformed header format")
	}
	if idx := bytes.Index(data, []byte("\r\n")); idx == 0 {
		return 0, true, nil
	}
	k, v := extractHeader(data)

	if !verifyHeaderKey(k) {
		fmt.Print("Error")
		return 0, false, errors.New("malformed header key value")
	}
	h[k] = v
	//return len(data) - 2, false, nil
	log.Println(h)
	return len(k) + len(v) + 4, false, nil
}

func haveCompleteData(data []byte) bool {
	idx := bytes.Index(data, []byte("\r\n"))
	return idx != -1
}

func verifyData(data []byte) bool {
	//"Host: localhost:42069\r\n\r\n"
	//data = bytes.TrimSpace(data) //Does not matter right?
	idx := bytes.Index(data, []byte(":"))
	return !unicode.IsSpace(rune(data[idx-1]))
}

func extractHeader(data []byte) (string, string) {
	idx := bytes.Index(data, []byte(":"))
	k := string(data[:idx])
	v := string(data[idx+1:])
	k, v = formatHeaders(k, v)
	//k = strings.TrimSpace(k)
	//v = strings.TrimSpace(v)
	log.Printf("Extracted key: %s\n", k)
	log.Printf("Extracted value: %s\n", v)
	return k, v
}

func formatHeaders(k, v string) (string, string) {
	k = strings.TrimSpace(k)
	k = strings.ToLower(k)
	v = strings.TrimSpace(v)

	return k, v
}

func verifyHeaderKey(k string) bool {
	for _, c := range k {
		if c > unicode.MaxASCII {
			log.Printf("Bad character %s", string(c))
			return false
		}
	}
	return true
}
