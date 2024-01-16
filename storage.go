package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	GetAccounts() ([]Account, error)
	CreateUser(*User) error
	GetUserByID(int) (*User, error)
	GetUserByUsername(string) (*User, error)
	CreateAccount(*Account) error
	UpdateAccount(*Account) error
	GetAccountByID(int) (*Account, error)
	DeleteAccountByID(int) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connectionString :=
		"user=postgres dbname=postgres password=admin sslmode=disable"

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Init() error {
	err := s.CreateUserTable()
	if err != nil {
		return err
	}
	err = s.CreateAccountTable()
	if err != nil {
		return err
	}

	// err := s.DeleteAccountTable()
	// if err != nil {
	// 	return err
	// }
	// err = s.DeleteUserTable()
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (s *PostgresStore) CreateUserTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS "user"(
		id SERIAL PRIMARY KEY,
		username VARCHAR(50),
		password VARCHAR(72),

		created_at TIMESTAMP,
		updated_at TIMESTAMP
	);
	`

	_, err := s.db.Exec(query)

	if err == nil {
		log.Println("user table created if not exists")
	}

	return err
}

func (s *PostgresStore) CreateAccountTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS account (
			id SERIAL PRIMARY KEY,
			user_id INT,
			first_name VARCHAR(50),
			last_name VARCHAR(50),
			number SERIAL,
			balance INT,
			created_at TIMESTAMP,
			updated_at TIMESTAMP,
			
			CONSTRAINT fk_user
				FOREIGN KEY (user_id)
					REFERENCES "user"(id)
		);
		`

	_, err := s.db.Exec(query)

	if err == nil {
		log.Println("account table created if not exists")
	}

	return err
}

func (s *PostgresStore) DeleteAccountTable() error {
	query := `DROP TABLE IF EXISTS account;`

	_, err := s.db.Exec(query)

	if err == nil {
		log.Println("account table deleted if exists")
	}

	return err
}

func (s *PostgresStore) DeleteUserTable() error {
	query := `DROP TABLE IF EXISTS "user";`

	_, err := s.db.Exec(query)

	if err == nil {
		log.Println("user table deleted if exists")
	}

	return err
}

func (s *PostgresStore) GetUserByUsername(username string) (*User, error) {
	query := `SELECT * FROM "user" WHERE username=$1`

	user := new(User)

	err := s.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *PostgresStore) GetUserByID(id int) (*User, error) {
	query := `SELECT * FROM "user" WHERE id=$1`

	user := new(User)

	err := s.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *PostgresStore) GetAccounts() ([]Account, error) {
	query := `SELECT * FROM account`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []Account{}
	for rows.Next() {
		a := Account{}

		err := rows.Scan(
			&a.ID,
			&a.FirstName,
			&a.LastName,
			&a.Number,
			&a.Balance,
			&a.CreatedAt,
			&a.UpdatedAt,
		)

		if err != nil {
			log.Println("failed to read from account row")
			return nil, err
		}

		fmt.Println("account: ", a)

		accounts = append(accounts, a)
	}

	return accounts, nil
}

func (s *PostgresStore) CreateUser(user *User) error {
	query := `
		INSERT INTO "user" (
			username, password, created_at, updated_at
		)
		VALUES($1, $2, $3, $4);
		`
	_, err := s.db.Query(
		query,
		user.Username,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) CreateAccount(a *Account) error {
	query := `
		INSERT INTO account (
			user_id, first_name, last_name, balance, number, created_at, updated_at
		)
		VALUES($1, $2, $3, $4, $5, $6, $7);
		`

	_, err := s.db.Query(
		query,
		a.UserID,
		a.FirstName,
		a.LastName,
		a.Balance,
		a.Number,
		a.CreatedAt,
		a.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	query := `SELECT * FROM account WHERE id=$1;`

	account := new(Account)

	err := s.db.QueryRow(query, id).Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get account by ID: %d", id)
	}

	return account, nil
}

func (s *PostgresStore) DeleteAccountByID(id int) error {
	query := `DELETE FROM account WHERE id=$1 RETURNING id`

	ID := new(int)

	err := s.db.QueryRow(query, id).Scan(ID)

	if err != nil {
		return errors.New("failed to delete account")
	}

	log.Printf("account with ID: %d deleted\n", id)

	return nil
}

func (s *PostgresStore) UpdateAccount(a *Account) error {
	return nil
}
