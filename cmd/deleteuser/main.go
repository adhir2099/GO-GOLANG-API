package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8098/user/delete/1", nil)
	if err != nil {
		log.Fatal(err)
	}

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
