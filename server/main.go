package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"unicode"
)

type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      int             `json:"id"`
}

type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   *Error      `json:"error"`
	ID      int         `json:"id"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func CountLetters(s string) int {
	count := 0
	for _, r := range s {
		if unicode.IsLetter(r) {
			count++
		}
	}
	return count
}

func CountWords(s string) int {
	return len(strings.Fields(s))
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		return
	}

	var result interface{}
	var err *Error

	switch req.Method {
	case "reverse":
		var params []string
		if err := json.Unmarshal(req.Params, &params); err != nil {
			log.Printf("Error decoding params: %v", err)
			return
		}
		if len(params) != 1 {
			err = &Error{Code: -32602, Message: "Invalid params"}
		} else {
			result = Reverse(params[0])
		}
	case "countLetters":
		var params []string
		if err := json.Unmarshal(req.Params, &params); err != nil {
			log.Printf("Error decoding params: %v", err)
			return
		}
		if len(params) != 1 {
			err = &Error{Code: -32602, Message: "Invalid params"}
		} else {
			result = CountLetters(params[0])
		}
	case "countWords":
		var params []string
		if err := json.Unmarshal(req.Params, &params); err != nil {
			log.Printf("Error decoding params: %v", err)
			return
		}
		if len(params) != 1 {
			err = &Error{Code: -32602, Message: "Invalid params"}
		} else {
			result = CountWords(params[0])
		}
	default:
		err = &Error{Code: -32601, Message: "Method not found"}
	}

	res := Response{
		JSONRPC: "2.0",
		Result:  result,
		Error:   err,
		ID:      req.ID,
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func main() {
	http.HandleFunc("/", HandleRequest)
	log.Println("Starting JSON-RPC server on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
