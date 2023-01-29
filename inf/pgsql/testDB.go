package pgsql

import (
	"errors"
	"userSL/models"
)

type testBD struct {
}

var (
	admin = models.User{Login: "admin", Password: "admin", Rule: models.Admin, Name: "first name Admin", LastName: "last name Admin", Dob: "01-01-1970"}
	user  = models.User{Login: "user", Password: "user", Rule: models.Read, Name: "first name user", LastName: "last name user", Dob: "21-10-2000"}
)

func GetTestDb() *testBD {
	return &testBD{}
}

func (*testBD) Load(login string) (models.User, error) {

	if login == "admin" {
		return admin, nil
	}
	if login == "user" {
		return user, nil
	}
	return models.User{}, errors.New("user not found")
}

func (*testBD) LoadAll() ([]models.User, error) {
	users := []models.User{admin, user}

	return users, nil
}

func (*testBD) Save(user *models.User) error {
	return nil
}
func (*testBD) Change(oldLogin string, user *models.User) error {
	if oldLogin == "admin" || oldLogin == "user" {
		return nil
	}
	return errors.New("user not found")
}
func (*testBD) Remove(login string, rule int) error {
	if login == "admin" || login == "user" {
		return nil
	}
	return errors.New("user not found")
}

func (*testBD) LastAdmin() bool {
	return true
}

func (*testBD) CloseDB() error {
	return nil
}
