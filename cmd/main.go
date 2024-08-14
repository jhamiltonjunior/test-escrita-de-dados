package main

import (
	"database/sql"
	"fmt"
	connectdb "github.com/jhamiltonjunior/test-escrita-de-dados/model"
)

func main() {
	mysql := connectdb.MySQLConnect{}
	db, err := connectdb.NewConnectionDB(&mysql, "127.0.0.1", "root", "0000", "teste_escrita_de_dados")
	if err != nil {
		return
	}

	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		panic(err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)

	// Iterando sobre os resultados
	for rows.Next() {
		var id int
		var name, email string
		err := rows.Scan(&id, &name, &email)
		if err != nil {
			panic(err)
		}
		fmt.Printf("ID: %d, Name: %s, Email: %s\n", id, name, email)
	}

	fmt.Println("cu")
}
