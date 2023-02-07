package db

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Connection struct {
	Host     string
	Port     string
	Password string
	User     string
	DbName   string
}

func (a *Connection) FromRequestBody(c echo.Context) Connection {
	u := Connection{}
	err := c.Bind(&u)
	if err != nil {
		panic(c.String(http.StatusBadRequest, "bad request"))
	}
	return u
}
func (a *Connection) ToString(ssl bool) string {
	var sslVal string
	if ssl {
		sslVal = "enable"
	} else {
		sslVal = "disable"
	}
	connTemplate := "user=%s password=%s dbname=%s host=%s port=%s sslmode=%s"
	connStr := fmt.Sprintf(connTemplate,
		a.User,
		a.Password,
		a.DbName,
		a.Host,
		a.Port,
		sslVal)
	return connStr
}

func (a Connection) GetConnection() (*sql.DB, error) {
	return sql.Open("postgres", a.ToString(false))
}
