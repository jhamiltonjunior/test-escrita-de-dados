package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	connectdb "github.com/jhamiltonjunior/test-escrita-de-dados/model"
	"io"
	"net/http"
	"time"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"username"`
}

func main() {
	initialTime := time.Now()
	mysql := connectdb.MySQLConnect{}

	actual := 0
	ln := 1000

	for i := 0; i < ln; i++ {
		users := HTTPGet("https://jsonplaceholder.typicode.com/users")

		inserirVariasLinhas(&mysql, users, ln*len(users), &actual)
	}

	migrarDados(&mysql)

	finalTime := time.Now()
	fmt.Printf("\nTempo de execucao do Script: %v", finalTime.Sub(initialTime))
}

func HTTPGet(url string) []User {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	//fmt.Println(resp.StatusCode)

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

func inserirVariasLinhas(mysql *connectdb.MySQLConnect, users []User, ln int, actual *int) {
	db, err := connectdb.NewConnectionDB(mysql, "127.0.0.1", "root", "0000", "teste_leitura_de_dados")
	if err != nil {
		panic(err)
	}
	defer func(close *sql.DB) {
		if err = close.Close(); err != nil {
			panic(err)
		}
	}(db)

	//_, err = db.Exec("DELETE FROM users;")
	//if err != nil {
	//	panic(err)
	//}

	for _, user := range users {
		var name, email, password = user.Name, user.Email, user.Password

		query := fmt.Sprintf(
			"INSERT INTO `users`(`name`, `email`, password) VALUES ('%s', '%s', '%s')",
			name, email, password,
		)

		_, err = db.Exec(query)
		if err != nil {
			panic(err)
		}
		*actual++
		fmt.Printf("\rSalvando %d/%d", *actual, ln)
	}
}

func migrarDados(mysql *connectdb.MySQLConnect) {
	db, err := connectdb.NewConnectionDB(mysql, "127.0.0.1", "root", "0000", "teste_leitura_de_dados")
	if err != nil {
		return
	}

	rows, err := db.Query("SELECT id, name, email, password FROM users;")
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

	//fmt.Println(User)

	dbBeWriter, err := connectdb.NewConnectionDB(mysql, "127.0.0.1", "root", "0000", "teste_escrita_de_dados")
	if err != nil {
		return
	}

	_, err = dbBeWriter.Exec("DELETE FROM users;")
	if err != nil {
		panic(err)
	}

	actual := 1

	// Iterando sobre os resultados
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
	fmt.Println("\nTransferência de dados concluída")
}
