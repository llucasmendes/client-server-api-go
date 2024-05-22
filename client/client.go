package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		logError(err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logError(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logError(fmt.Errorf("unexpected status code: %d", resp.StatusCode))
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logError(err)
		return
	}

	err = ioutil.WriteFile("cotacao.txt", []byte(fmt.Sprintf("Dólar: %s", body)), 0644)
	if err != nil {
		logError(err)
	}
}

func logError(err error) {
	fmt.Println("Error:", err)
}
