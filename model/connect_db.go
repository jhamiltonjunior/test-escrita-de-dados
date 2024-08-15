package connect_db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type ConnectDB interface {
	Connection() (*sql.DB, error)
}

type MySQLConnect struct {
	host, user, pass, database string
}

func (m *MySQLConnect) Connection() (*sql.DB, error) {
	dsn := m.user + ":" + m.pass + "@tcp(" + m.host + ":3306)/" + m.database
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

func NewConnectionDB(host, user, pass, db string) ConnectDB {
	return &MySQLConnect{
		host:     host,
		pass:     pass,
		user:     user,
		database: db,
	}
}
