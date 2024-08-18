package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	connectdb "github.com/jhamiltonjunior/test-escrita-de-dados/model"
	"io"
	"net/http"
	"runtime"
	"time"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"username"`
}

func printStatus(totalInserts int, start time.Time) {
	// Limpa a tela (opcional, dependendo do sistema operacional)
	fmt.Print("\033[H\033[2J")

	// Tempo de execução
	elapsed := time.Since(start)

	// Memória utilizada
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	alloc := m.Alloc / 1024 / 1024 // Memória alocada em MB

	// Imprimindo as informações
	fmt.Printf("Total de Dados Inseridos: %d\n", totalInserts)
	fmt.Printf("Tempo de Execução: %s\n", elapsed)
	fmt.Printf("Memória Utilizada: %d MB\n", alloc)
}

func main() {
	done := make(chan bool, 10)
	initialTime := time.Now()

	connReader := connectdb.NewConnectionDB("127.0.0.1", "root", "0000", "teste_leitura_de_dados")

	dbReader, err := connReader.Connection()
	if err != nil {
		panic(err)
	}
	defer func(close *sql.DB) {
		if err = close.Close(); err != nil {
			panic(err)
		}
	}(dbReader)

	connWriter := connectdb.NewConnectionDB("127.0.0.1", "root", "0000", "teste_escrita_de_dados")

	dbWriter, err := connWriter.Connection()
	if err != nil {
		panic(err)
	}
	defer func(close *sql.DB) {
		if err = close.Close(); err != nil {
			panic(err)
		}
	}(dbWriter)

	start := time.Now()

	totalInserts := 0
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			totalInserts += 1

			printStatus(totalInserts, start)
		}
	}()

	actual := 0
	ln := 1000

	for i := 0; i < ln; i++ {
		users := HTTPGet("https://jsonplaceholder.typicode.com/users")

		inserirVariasLinhas(dbReader, users, ln*len(users), &actual, done)
	}

	migrarDados(dbReader, dbWriter, done)

	finalTime := time.Now()
	fmt.Printf("\nTempo de execucao do Script: %v", finalTime.Sub(initialTime))
}

func HTTPGet(url string) []User {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	body, _ := io.ReadAll(resp.Body)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	var users []User

	err = json.Unmarshal(body, &users)
	if err != nil {
		panic(err)
	}

	return users
}

func inserirVariasLinhas(db *sql.DB, users []User, ln int, actual *int, done chan bool) {

	for _, user := range users {
		var name, email, password = user.Name, user.Email, user.Password

		query := fmt.Sprintf(
			"INSERT INTO `users`(`name`, `email`, password) VALUES ('%s', '%s', '%s')",
			name, email, password,
		)

		_, err := db.Exec(query)
		if err != nil {
			panic(err)
		}
		*actual++
		fmt.Printf("\rSalvando %d/%d", *actual, ln)
	}
}

func migrarDados(dbReader, dbBeWriter *sql.DB, done chan bool) {
	rows, err := dbReader.Query("SELECT id, name, email, password FROM users;")
	if err != nil {
		panic(err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)

	totalLinhas := 0

	var users []User

	for rows.Next() {

		totalLinhas++

		var id int
		var name, email, password string
		err := rows.Scan(&id, &name, &email, &password)
		if err != nil {
			panic(err)
		}

		users = append(users, User{id, name, email, password})
	}

	actual := 1

	for _, user := range users {

		var name, email, password = user.Name, user.Email, user.Password

		query := fmt.Sprintf(
			"INSERT INTO `users`(`name`, `email`, password) VALUES ('%s', '%s', '%s')",
			name, email, password,
		)

		_, err = dbBeWriter.Exec(query)
		if err != nil {
			panic(err)
		}
		fmt.Printf("\rSalvando %d/%d", actual, totalLinhas)
		actual++
	}
	done <- true
	fmt.Println("\nTransferência de dados concluída")

	close(done)
}
