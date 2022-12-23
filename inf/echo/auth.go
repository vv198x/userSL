package echo

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"userSL/inf/pgsql"
	"userSL/models"
	"userSL/pkg/config"
)

// getToken godoc
// @Tags auth
// @Summary Authentication
// @Produce json
// @Param message body models.User true  "New user"
// @Success 200 {object} JSONLogin{login=string,password=string}
// @Failure	400 {object} JSONResult{message=string} "Bad Request"
// @Failure	401 {object} JSONResult{message=string} "Authentication error"
// @Failure	500
// @Router / [get]
func getToken(c echo.Context) error {
	res := new(JSONLogin)

	err := c.Bind(res)
	err = c.Validate(res)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	db, _ := c.Get("db").(pgsql.Storage)
	user, err := db.Load(res.Login)

	if res.Password != user.Password {
		log.Println("Wrong password ")
		return echo.NewHTTPError(http.StatusUnauthorized, "Authentication error")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": user.Login,
		"rule":  user.Rule,
	})

	tokenString, err := token.SignedString([]byte(*config.JWTkey))
	if err != nil {
		log.Println("Can't create token ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't create token")
	}

	return c.JSON(http.StatusOK, JSONToken{tokenString})
}

func checkToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tVal := c.FormValue("token")

		if tVal == "" {
			log.Println("Token empty ")
			return echo.NewHTTPError(http.StatusUnauthorized, "Authentication error")
		}

		token, err := jwt.Parse(tVal, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected Signing Method: %v", token.Header["alg"])
			}

			return []byte(*config.JWTkey), nil
		})

		if !token.Valid || err != nil {

			return echo.NewHTTPError(http.StatusUnauthorized, "Authentication error")
		}
		if token.Claims.(jwt.MapClaims)["login"] != "" {

			c.Set("rule", token.Claims.(jwt.MapClaims)["rule"])
		}

		return next(c)
	}
}

func forAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		rule := c.Get("rule").(float64)
		if rule != models.Admin {
			return echo.NewHTTPError(http.StatusForbidden, "Forbidden")
		}
		return next(c)
	}
}

func forAll(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		rule := c.Get("rule").(float64)
		if rule == models.Lock {
			return echo.NewHTTPError(http.StatusLocked, "Locked")
		}
		return next(c)
	}
}
