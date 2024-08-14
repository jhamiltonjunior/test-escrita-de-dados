package connect_db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type ConnectDB interface {
	Connection(host, user, pass, db string) (*sql.DB, error)
}

type MySQLConnect struct{}

func (m *MySQLConnect) Connection(host, user, pass, database string) (*sql.DB, error) {
	dsn := user + ":" + pass + "@tcp(" + host + ":3306)/" + database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco de dados: %w", err)
	}

	// Verifica se a conexão está ativa
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("não foi possível pingar o banco de dados: %w", err)
	}

	return db, nil
}

func NewConnectionDB(c ConnectDB, host, user, pass, db string) (*sql.DB, error) {
	return c.Connection(host, user, pass, db)
}
