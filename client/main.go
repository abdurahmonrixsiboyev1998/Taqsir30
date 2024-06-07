package main

import (
    "bufio"
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "strings"
)

type Request struct {
    JSONRPC string        `json:"jsonrpc"`
    Method  string        `json:"method"`
    Params  []interface{} `json:"params"`
    ID      int           `json:"id"`
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

func sendRequest(method string, params []interface{}) (*Response, error) {
    req := Request{
        JSONRPC: "2.0",
        Method:  method,
        Params:  params,
        ID:      1,
    }
    reqBody, err := json.Marshal(req)
    if err != nil {
        return nil, err
    }

    resp, err := http.Post("http://localhost:8080", "application/json", bytes.NewBuffer(reqBody))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var res Response
    if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
        return nil, err
    }

    return &res, nil
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Print("Enter a string: ")
        input, _ := reader.ReadString('\n')
        input = strings.TrimSpace(input)

        if input == "exit" {
            break
        }

        res, err := sendRequest("reverse", []interface{}{input})
        if err != nil {
            log.Fatalf("Error sending request: %v", err)
        }
        fmt.Printf("Reversed string: %v\n", res.Result)

        res, err = sendRequest("countLetters", []interface{}{input})
        if err != nil {
            log.Fatalf("Error sending request: %v", err)
        }
        fmt.Printf("Number of letters: %v\n", res.Result)

        res, err = sendRequest("countWords", []interface{}{input})
        if err != nil {
            log.Fatalf("Error sending request: %v", err)
        }
        fmt.Printf("Number of words: %v\n", res.Result)
    }
}
