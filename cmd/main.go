package main

import (
	"database/sql"
	"fmt"
	connectdb "github.com/jhamiltonjunior/test-escrita-de-dados/model"
)

func main() {
	mysql := connectdb.MySQLConnect{}
	db, err := connectdb.NewConnectionDB(&mysql, "127.0.0.1", "root", "0000", "teste_leitura_de_dados")
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

	db2, err := connectdb.NewConnectionDB(&mysql, "127.0.0.1", "root", "0000", "teste_escrita_de_dados")
	if err != nil {
		return
	}

	password := "securepassword"
	createAt := "2024-08-14 10:00:00"
	updateAt := "2024-08-14 10:00:00"
	deleteAt := "NULL" // ou use "2024-08-14 10:00:00" para uma data
	active := 1

	// Iterando sobre os resultados
	for rows.Next() {
		var id int
		var name, email string
		err := rows.Scan(&id, &name, &email)
		if err != nil {
			panic(err)
		}
		fmt.Printf("ID: %d, Name: %s, Email: %s\n", id, name, email)

		query := fmt.Sprintf(
			"INSERT INTO `users`(`id`, `name`, `email`, `password`, `create_at`, `update_at`, `delete_at`, `active`) VALUES ('%d', '%s', '%s', '%s', '%s', '%s', '%s', '%d')",
			id, name, email, password, createAt, updateAt, deleteAt, active,
		)

		_, err = db2.Query(query)
		if err != nil {
			return
		}
	}

	fmt.Println("cu")
}
