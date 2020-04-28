package users

import (
	"bookstore_users-api/datasources/mysql/users"
	"bookstore_users-api/utils/dates"
	"bookstore_users-api/utils/errors"
	"fmt"
	"strings"
)

const (
	queryInsertUser  = "INSERT INTO users(first_name, last_name, email, date_created) VALUES (?, ?, ?, ?);"
	indexUniqueEmail = "email_uindex"
)

var (
	usersDB = make(map[int64]*User)
)

func (user *User) Get(userId int64) *errors.RestErr {
	result := usersDB[userId]
	if result == nil {
		return errors.NewNotFoundError(fmt.Sprintf("user %d not found", userId))
	}

	user.Id = result.Id
	user.FirstName = result.FirstName
	user.LastName = result.LastName
	user.Email = result.Email
	user.DateCreated = result.DateCreated

	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	user.DateCreated = dates.GetNowString()

	insertResult, err := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated)
	if err != nil {
		if strings.Contains(err.Error(), indexUniqueEmail) {
			return errors.NewBadRequestError(fmt.Sprintf("email %s already exists", user.Email))
		}
		return errors.NewInternalServerError(fmt.Sprintf("error on trying to save a user: %s", err.Error()))
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("error on trying to save a user: %s", err.Error()))
	}

	user.Id = userId
	return nil
}
