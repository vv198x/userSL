package echo

import (
	"github.com/labstack/echo"
	"log"
	"net/http"
	"strings"
	"userSL/inf/pgsql"
	"userSL/models"
)

// Read godoc
// @Tags read
// @Summary Retrieves user based on given Login
// @Produce json
// @Param login path string true "User login"
// @Success 200 {object} models.User
// @Failure	404 {object} models.JSONResult{message=string} "Not found"
// @Failure	500
// @Router /{login} [get]
func Read(c echo.Context) error {
	user := *(c.Get("user").(*models.User))

	//Не выводить пароль. Так занулил, если тип сменится.
	user.Password = models.User{}.Password
	return c.JSON(http.StatusOK, &user)
}

// Read godoc
// @Tags read
// @Summary Retrieves users
// @Produce json
// @Success 200 {object} models.User
// @Failure	500
// @Router / [get]
func ReadAll(c echo.Context) error {
	db, _ := c.Get("db").(pgsql.Storage)
	users, err := db.LoadAll()
	if err == nil {
		for i := range users {
			users[i].Password = models.User{}.Password
		}
		return c.JSON(http.StatusOK, &users)
	}

	log.Println("DB error ", err)
	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
}

// Create godoc
// @Summary Create new user
// @Tags admins
// @Produce json
// @Param message body models.User true  "New user"
// @Success 201 {object} models.User
// @Failure	400 {object} models.JSONResult{message=string} "Validation error"
// @Failure	409 {object} models.JSONResult{message=string} "User with this login exists"
// @Failure	500
// @Router / [post]
func Create(c echo.Context) error {
	user := c.Get("validUser").(*models.User)
	db, _ := c.Get("db").(pgsql.Storage)

	err := db.Save(user)
	if err == nil {
		return c.JSON(http.StatusCreated, user)
	}
	//При OnConflict возращается "pg: no rows in result set"
	if strings.Contains(err.Error(), "no rows") {
		return echo.NewHTTPError(http.StatusConflict, "User exists")
	}

	log.Println("DB error ", err)
	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
}

// Update godoc
// @Summary Update user on given Login
// @Tags admins
// @Produce json
// @Param login path string true "User login"
// @Param message body models.User true  "Update user"
// @Success 202 {object} models.User
// @Failure	404 {object} models.JSONResult{message=string} "Not found"
// @Failure	400 {object} models.JSONResult{message=string} "Validation error"
// @Failure	500
// @Router /{login} [put]
func Update(c echo.Context) error {
	db, _ := c.Get("db").(pgsql.Storage)
	user := c.Get("validUser").(*models.User)

	oldLogin := c.Param("login")

	// Передаю ссылку на user, для sql.returning
	err := db.Change(oldLogin, user)
	if err == nil {
		return c.JSON(http.StatusAccepted, user)
	}

	log.Println("DB error ", err)
	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
}

// Delete godoc
// @Tags admins
// @Summary Delete user on given Login
// @Produce json
// @Param login path string true "User login"
// @Success 200
// @Failure	400 {object} models.JSONResult{message=string} "Attempt to remove the last admin"
// @Failure	404 {object} models.JSONResult{message=string} "Not found"
// @Failure	500
// @Router /{login} [delete]
func Delete(c echo.Context) error {
	db, _ := c.Get("db").(pgsql.Storage)
	user := *(c.Get("user").(*models.User))

	err := db.Remove(user.Login, user.Rule)

	if err == nil {
		return c.NoContent(http.StatusOK)
	}

	log.Println("DB error ", err)
	return echo.NewHTTPError(http.StatusOK, err.Error())
}
