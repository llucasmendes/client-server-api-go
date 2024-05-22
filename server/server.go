package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	apiURL      = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	dbTimeout   = 10 * time.Millisecond
	apiTimeout  = 500 * time.Millisecond // Aumentado para 500ms
	serverPort  = ":8080"
	createTable = `CREATE TABLE IF NOT EXISTS cotacoes (id INTEGER PRIMARY KEY, bid TEXT, timestamp DATETIME DEFAULT CURRENT_TIMESTAMP)`
	insertQuote = `INSERT INTO cotacoes (bid) VALUES (?)`
)

type QuoteResponse struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

func main() {
	http.HandleFunc("/cotacao", handleCotacao)
	fmt.Println("Server running on port", serverPort)
	http.ListenAndServe(serverPort, nil)
}

func handleCotacao(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), apiTimeout)
	defer cancel()

	quote, err := fetchQuote(ctx)
	if err != nil {
		http.Error(w, "Failed to fetch quote", http.StatusInternalServerError)
		logError(err)
		return
	}

	ctxDB, cancelDB := context.WithTimeout(context.Background(), dbTimeout)
	defer cancelDB()

	err = saveQuote(ctxDB, quote)
	if err != nil {
		http.Error(w, "Failed to save quote", http.StatusInternalServerError)
		logError(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(quote))
}

func fetchQuote(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var quoteResponse QuoteResponse
	err = json.NewDecoder(resp.Body).Decode(&quoteResponse)
	if err != nil {
		return "", err
	}

	return quoteResponse.USDBRL.Bid, nil
}

func saveQuote(ctx context.Context, bid string) error {
	db, err := sql.Open("sqlite3", "./cotacoes.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.ExecContext(ctx, createTable)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, insertQuote, bid)
	return err
}

func logError(err error) {
	fmt.Println("Error:", err)
}
