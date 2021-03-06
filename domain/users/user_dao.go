package users

import (
	"fmt"
	"github.com/n0rdy/bookstore_users-api/datasources/mysql/users"
	"github.com/n0rdy/bookstore_users-api/logger"
	"github.com/n0rdy/bookstore_users-api/utils/mysql"
	"github.com/n0rdy/bookstore_utils-go/rest_errors"
	"strings"
)

const (
	queryInsertUser                 = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES (?, ?, ?, ?, ?, ?);"
	queryGetUser                    = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id = ?;"
	queryUpdateUser                 = "UPDATE users SET first_name = ?, last_name = ?, email = ? WHERE id = ?;"
	queryDeleteUser                 = "DELETE FROM users WHERE id = ?;"
	queryFindUserByStatus           = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status = ?;"
	queryFindUserByEmailAndPassword = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE email = ? AND password = ? ANS status = ?;"
)

func (user *User) Get() rest_errors.RestErr {
	stmt, err := users.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error on trying to prepare get user statement", err)
		return rest_errors.NewInternalServerError("database error", err)
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		logger.Error("error on trying to get user by id", getErr)
		return rest_errors.NewInternalServerError("database error", getErr)
	}

	return nil
}

func (user *User) Save() rest_errors.RestErr {
	stmt, err := users.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error on trying to prepare save user statement", err)
		return rest_errors.NewInternalServerError("database error", err)
	}
	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if saveErr != nil {
		logger.Error("error on trying to save user", saveErr)
		return rest_errors.NewInternalServerError("database error", saveErr)
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error on trying to get last insert id after creating a new user", err)
		return rest_errors.NewInternalServerError("database error", err)
	}

	user.Id = userId
	return nil
}

func (user *User) Update() rest_errors.RestErr {
	stmt, err := users.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error on trying to prepare update user statement", err)
		return rest_errors.NewInternalServerError("database error", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		logger.Error("error on trying to update user", err)
		return rest_errors.NewInternalServerError("database error", err)
	}

	return nil
}

func (user *User) Delete() rest_errors.RestErr {
	stmt, err := users.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error on trying to prepare delete user statement", err)
		return rest_errors.NewInternalServerError("database error", err)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.Id); err != nil {
		logger.Error("error on trying to delete user", err)
		return rest_errors.NewInternalServerError("database error", err)
	}
	return nil
}

func (user *User) FindByStatus(status string) ([]User, rest_errors.RestErr) {
	stmt, err := users.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		logger.Error("error on trying to prepare find user by status statement", err)
		return nil, rest_errors.NewInternalServerError("database error", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error on trying to find user by status", err)
		return nil, rest_errors.NewInternalServerError("database error", err)
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			logger.Error("error on trying to scan user raw into user struct", err)
			return nil, rest_errors.NewInternalServerError("database error", err)
		}
		results = append(results, user)
	}

	if len(results) == 0 {
		return nil, rest_errors.NewNotFoundError(fmt.Sprintf("No user matches status %s", status))
	}
	return results, nil
}

func (user *User) FindByEmailAndPassword() rest_errors.RestErr {
	stmt, err := users.Client.Prepare(queryFindUserByEmailAndPassword)
	if err != nil {
		logger.Error("error on trying to prepare find user by email and password statement", err)
		return rest_errors.NewInternalServerError("database error", err)
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Email, user.Password, StatusActive)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		if strings.Contains(getErr.Error(), mysql.ErrorNoRows) {
			return rest_errors.NewNotFoundError("Invalid user credentials")
		}
		logger.Error("error on trying to get user by email and password", getErr)
		return rest_errors.NewInternalServerError("database error", getErr)
	}

	return nil
}
