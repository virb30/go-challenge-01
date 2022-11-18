package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type CotacaoDto struct {
	Usdbrl struct {
		Valor string `json:"bid"`
	} `json:"USDBRL"`
}

type Quotation struct {
	ID    int `gorm:"primaryKey"`
	Valor string
}

func main() {
	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	dsn := "database.sqlite3"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&Quotation{})
	if err != nil {
		log.Fatal(err)
	}
	c, err := getCotacao()
	if err != nil {
		log.Fatal(err)
	}
	cotacao := Quotation{Valor: c.Usdbrl.Valor}
	insertCotacao(db, &cotacao)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"valor": %s}`, c.Usdbrl.Valor)))
}

func getCotacao() (*CotacaoDto, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(200*time.Millisecond))
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
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
	var c CotacaoDto
	err = json.Unmarshal(body, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func insertCotacao(db *gorm.DB, c *Quotation) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Millisecond))
	defer cancel()
	db.WithContext(ctx).Create(&c)
}
