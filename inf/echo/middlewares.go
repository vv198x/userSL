package echo

import (
	"github.com/labstack/echo"
	"log"
	"net/http"
	"userSL/inf/pgsql"
	"userSL/models"
)

func forAdmin(login, pass string, c echo.Context) (bool, error) {
	db, _ := c.Get("db").(pgsql.Storage)
	u, err := db.Load(login)
	if err == nil && u.Password == pass {

		if u.Rule == models.Admin {
			return true, nil
		} else {
			log.Println("No access rights ", login)
			c.String(http.StatusForbidden, "forbidden")
		}
	}
	return false, err
}

func forAll(login, pass string, c echo.Context) (bool, error) {
	db, _ := c.Get("db").(pgsql.Storage)
	u, err := db.Load(login)
	if err == nil && u.Password == pass {

		if u.Rule == models.Lock {
			log.Println("No access rights ", login)
			c.String(http.StatusForbidden, "forbidden")
			return false, err
		} else {
			return true, nil
		}
	}
	return false, err
}
