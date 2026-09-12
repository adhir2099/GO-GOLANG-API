package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	req, err := http.NewRequest(http.MethodPatch, "http://localhost:8098/user/update/1", bytes.NewBufferString(`{"name":"User1","email":"user1@user.com"}`))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n%s\n", response.Status, body)
}
