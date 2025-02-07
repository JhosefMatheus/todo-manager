package test_utils

import (
	"database/sql"
	"testing"
	"todo-manager/models"

	"github.com/joho/godotenv"
)

func SetupEnv(t *testing.T, envPath string) {
	if err := godotenv.Load(envPath); err != nil {
		t.Errorf("Erro ao carregar .env: %v", err)
	}
}

func SetupUserTable(db *sql.DB, t *testing.T) {
	sql := `
		insert into user (name, email, password)
		value ('Jhosef Matheus', 'jhosef.dev@gmail.com', sha2('9=0=y7MA5S>y', 256));
	`

	if _, err := db.Exec(sql); err != nil {
		t.Errorf("Erro ao inserir usuário: %v", err)
	}
}

func ClearUserTable(db *sql.DB) {
	sql := `
		delete from user;
	`

	db.Exec(sql)
}

func GetInsertedUser(db *sql.DB, t *testing.T) models.UserModel {
	var insertedUser models.UserModel

	sql := `
		select
			id,
			name,
			email,
			created_at,
			updated_at
		from user
		where email = ?
		limit 1;
	`

	rows, err := db.Query(sql, "jhosef.dev@gmail.com")

	if err != nil {
		t.Error(err)
	}

	if rows.Next() {
		if err = rows.Scan(&insertedUser.Id, &insertedUser.Name, &insertedUser.Email, &insertedUser.CreatedAt, &insertedUser.UpdatedAt); err != nil {
			t.Errorf("Erro ao ler resposta do banco de dados: %v", err)
		}
	} else {
		t.Error("Nenhum usuário encontrado")
	}

	return insertedUser
}

func SetupProjectTable(db *sql.DB, t *testing.T) int64 {
	sql := `
		insert into project (name)
		value ('Teste');
	`

	result, err := db.Exec(sql)

	if err != nil {
		t.Errorf("Erro ao inserir projeto: %v", err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		t.Errorf("Erro ao pegar id do projeto inserido: %v", err)
	}

	return id
}

func ClearProjectTable(db *sql.DB) {
	sql := `
			delete from project;
		`

	db.Exec(sql)
}
