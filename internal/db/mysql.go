package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/google/uuid"
	"github.com/knudsenTaunus/employeeService/internal/config"
	"github.com/knudsenTaunus/employeeService/internal/model"
)

var mysqlErr *mysql.MySQLError

type Cryptor3000 interface {
	Encrypt(text string) (string, error)
	Decrypt(text string) (string, error)
}

type MySQL struct {
	conn    *sql.DB
	cryptor Cryptor3000
}

func NewMySQL(config *config.Config, cryptor Cryptor3000) (*MySQL, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.Mysqldatabase.User, config.Mysqldatabase.Password, config.Mysqldatabase.Host, config.Mysqldatabase.Port, "users")
	db, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, err
	}

	return &MySQL{conn: db, cryptor: cryptor}, nil
}

func (m *MySQL) GetAll() ([]model.User, error) {
	employees := make([]model.User, 0)
	stmt, err := m.conn.Prepare("SELECT id, first_name, last_name, nickname, password, email, country, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		tmp := model.User{}
		encryptedPassword := ""
		err := rows.Scan(&tmp.ID, &tmp.FirstName, &tmp.LastName, &tmp.Nickname, &encryptedPassword, &tmp.Email, &tmp.Country, &tmp.CreatedAt, &tmp.UpdatedAt)
		if err != nil {
			return nil, err
		}

		decryptedPassword, err := m.cryptor.Decrypt(encryptedPassword)
		if err != nil {
			return nil, err
		}
		tmp.Password = decryptedPassword

		employees = append(employees, tmp)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return employees, nil
}

func (m *MySQL) GetPaginatedAndFiltered(page, pageSize int, filter string) ([]model.User, error) {
	limit := page * pageSize
	offset := (page * pageSize) - pageSize
	filterKeyValue := strings.Split(filter, "=")

	queryStatement := fmt.Sprintf("SELECT * FROM users WHERE %s = '%s' LIMIT %d, %d", filterKeyValue[0], filterKeyValue[1], offset, limit)

	rows, err := m.conn.Query(queryStatement)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]model.User, 0, limit)
	for rows.Next() {
		tmp := model.User{}
		encryptedPassword := ""
		err := rows.Scan(&tmp.ID, &tmp.FirstName, &tmp.LastName, &tmp.Nickname, &encryptedPassword, &tmp.Email, &tmp.Country, &tmp.CreatedAt, &tmp.UpdatedAt)
		if err != nil {
			return nil, err
		}

		decryptedPassword, err := m.cryptor.Decrypt(encryptedPassword)
		if err != nil {
			return nil, err
		}
		tmp.Password = decryptedPassword
		users = append(users, tmp)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (m *MySQL) Get(id string) (model.User, error) {
	tmp := model.User{}
	stmt, err := m.conn.Prepare("SELECT id, first_name, last_name, nickname, password, email, country, created_at, updated_at FROM users WHERE id = ?")
	if err != nil {
		return model.User{}, err
	}

	var encryptedPassword string

	err = stmt.QueryRow(id).Scan(&tmp.ID, &tmp.FirstName, &tmp.LastName, &tmp.Nickname, &encryptedPassword, &tmp.Email, &tmp.Country, &tmp.CreatedAt, &tmp.UpdatedAt)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return tmp, model.NotFoundError
		default:
			return tmp, err
		}
	}

	decryptedPassword, err := m.cryptor.Decrypt(encryptedPassword)
	if err != nil {
		return model.User{}, err
	}

	tmp.Password = decryptedPassword

	return tmp, nil
}
func (m *MySQL) Create(u model.User) (model.User, error) {
	stmt, err := m.conn.Prepare("INSERT INTO users (id, first_name, last_name, nickname, password, email, country, created_at, updated_at) VALUES (?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return model.User{}, err
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return model.User{}, err
	}

	encryptedPassword, cryptErr := m.cryptor.Encrypt(u.Password)
	if cryptErr != nil {
		return model.User{}, cryptErr
	}

	_, err = stmt.Exec(id.String(), u.FirstName, u.LastName, u.Nickname, encryptedPassword, u.Email, u.Country, time.Now().UTC(), time.Now().UTC())
	if err != nil {
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			if strings.Contains(err.Error(), "users.email") {
				return model.User{}, model.DuplicateMailError
			}

			if strings.Contains(err.Error(), "users.nickname") {
				return model.User{}, model.DuplicateNickError
			}
		}
		return model.User{}, err
	}

	createdUser, err := m.Get(id.String())
	if err != nil {
		return model.User{}, err
	}

	return createdUser, nil
}

func (m *MySQL) Update(u model.User) (model.User, error) {
	stmt, err := m.conn.Prepare("UPDATE users SET first_name= ?, last_name= ?, password= ?, email=?, country=?, updated_at=? WHERE id= ?")
	if err != nil {
		return model.User{}, err
	}

	encryptedPassword, cryptErr := m.cryptor.Encrypt(u.Password)
	if cryptErr != nil {
		return model.User{}, cryptErr
	}

	_, err = stmt.Exec(u.FirstName, u.LastName, encryptedPassword, u.Email, u.Country, time.Now().UTC(), u.ID)
	if err != nil {
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			if strings.Contains(err.Error(), "users.email") {
				return model.User{}, model.DuplicateMailError
			}

			if strings.Contains(err.Error(), "users.nickname") {
				return model.User{}, model.DuplicateNickError
			}
		}
	}

	updatedUser, err := m.Get(u.ID)
	if err != nil {
		return model.User{}, err
	}

	return updatedUser, nil
}

func (m *MySQL) Delete(id string) error {
	stmt, err := m.conn.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
