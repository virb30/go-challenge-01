package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
	Valor float64 `json:"valor"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(300*time.Millisecond))
	defer cancel()
	c, err := getCotacao(ctx)
	if err != nil {
		log.Println(err)
	}
	err = saveCotacao(c, "cotacao.txt")
	if err != nil {
		log.Println(err)
	}
}

func getCotacao(ctx context.Context) (*Cotacao, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var c Cotacao
	json.Unmarshal(body, &c)
	return &c, nil
}

func saveCotacao(c *Cotacao, name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	f.Write([]byte(fmt.Sprintf("Dólar: %v", c.Valor)))
	return nil
}
