package users

import (
	"bookstore_users-api/datasources/mysql/users"
	"bookstore_users-api/utils/dates"
	"bookstore_users-api/utils/errors"
	"bookstore_users-api/utils/mysql"
)

const (
	errorNoRows     = "no rows in result set"
	queryInsertUser = "INSERT INTO users(first_name, last_name, email, date_created) VALUES (?, ?, ?, ?);"
	queryGetUser    = "SELECT id, first_name, last_name, email, date_created FROM users WHERE id = ?;"
)

func (user *User) Get(userId int64) *errors.RestErr {
	stmt, err := users.Client.Prepare(queryGetUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	result := stmt.QueryRow(userId)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated); getErr != nil {
		mysql.ParseError(getErr)
	}

	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	user.DateCreated = dates.GetNowString()

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated)
	if saveErr != nil {
		return mysql.ParseError(saveErr)
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		return mysql.ParseError(err)
	}

	user.Id = userId
	return nil
}
