package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	connectdb "github.com/jhamiltonjunior/test-escrita-de-dados/model"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"` // Corrigi o campo para "password"
}

func printStatus(totalInserts int, start time.Time) {
	// Limpa a tela (opcional, dependendo do sistema operacional)
	//fmt.Print("\033[H\033[2J")
	//
	//// Tempo de execução
	//elapsed := time.Since(start)
	//
	//// Memória utilizada
	//var m runtime.MemStats
	//runtime.ReadMemStats(&m)
	//alloc := m.Alloc / 1024 / 1024 // Memória alocada em MB
	//
	//// Imprimindo as informações
	////fmt.Printf("Total de Dados Inseridos: %d\n", totalInserts)
	////fmt.Printf("Tempo de Execução: %s\n", elapsed)
	////fmt.Printf("Memória Utilizada: %d MB\n", alloc)
}

func RunMigration() {
	done := make(chan struct{})
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

	//defer close(done)

	start := time.Now()

	totalInserts := 0
	go func() {
		for _ = range done {
			totalInserts++

			printStatus(totalInserts, start)
		}
	}()

	//actual := 0
	//ln := 1000
	//
	////for i := 0; i < ln; i++ {
	////	done <- struct{}{}
	////	users := HTTPGet("https://jsonplaceholder.typicode.com/users")
	////
	////	inserirVariasLinhas(dbReader, users, ln*len(users), &actual)
	////}

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

func inserirVariasLinhas(db *sql.DB, users []User, ln int, actual *int) {

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

func migrarDados(dbReader, dbBeWriter *sql.DB, done chan struct{}) {
	rows, err := dbReader.Query("SELECT id, name, email, password FROM users limit 2;")
	if err != nil {
		panic(err)
	}

	totalLinhas := 0

	var users chan User

	actual := 1

	go func() {
		for user := range users {

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
			done <- struct{}{}
		}

		close(users)
		close(done)
	}()

	for rows.Next() {

		totalLinhas++

		var id int
		var name, email, password string
		err = rows.Scan(&id, &name, &email, &password)
		if err != nil {
			panic(err)
		}

		fmt.Println(users)
		fmt.Println("users")
		users <- User{id, name, email, password}
	}

	err = rows.Close()
	if err != nil {
		panic(err)
	}

	fmt.Println("\nTransferência de dados concluída")
}

func main() {
	// Criar arquivo para salvar o perfil de CPU
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	// Iniciar perfil de CPU
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	memFile, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal("could not create MEM profile: ", err)
	}
	defer func(f *os.File) {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}(memFile)

	runtime.GC()

	if err = pprof.WriteHeapProfile(memFile); err != nil {
		log.Fatal("WriteHeapProfile: ", err)
	}

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	RunMigration()
}
