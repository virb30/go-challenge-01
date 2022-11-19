package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type CotacaoDto struct {
	Usdbrl struct {
		Valor string `json:"bid"`
	} `json:"USDBRL"`
}

type Cotacao struct {
	ID    int `gorm:"primaryKey"`
	Valor string
}

func main() {
	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	dsn := "database.sqlite3"
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Println(err)
	}
	err = migrate(db)
	if err != nil {
		log.Println(err)
	}
	c, err := getCotacao()
	if err != nil {
		log.Println(err)
	}
	cotacao := Cotacao{Valor: c.Usdbrl.Valor}
	err = insertCotacao(db, cotacao)
	if err != nil {
		log.Println(err)
	}
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

func insertCotacao(db *sql.DB, c Cotacao) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Millisecond))
	defer cancel()
	stmt, err := db.PrepareContext(ctx, "INSERT INTO cotacoes (valor) values (?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, c.Valor)
	if err != nil {
		return err
	}
	return nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS cotacoes (id integer primary key autoincrement, valor varchar)")
	if err != nil {
		return err
	}
	return nil
}
