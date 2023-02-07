package user

import (
	"fmt"
	"gotest/db"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func GetInfo(c echo.Context) error {
	dbSess := c.Get("database").(*sessions.Session)
	return c.JSON(http.StatusOK, fmt.Sprintf("You're connected as %s", dbSess.Values["db"].(db.Connection).User))
}

func Connect(c echo.Context) error {
	con := (&db.Connection{}).FromRequestBody(c)
	uDb, err0 := con.GetConnection()
	pingError := uDb.Ping()
	if pingError != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("unable to connect using credentials provided %d", pingError)})
	}
	if err0 != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "could not connect"})
	}
	defer uDb.Close()
	dbSess := c.Get("database").(*sessions.Session)
	dbSess.Values["db"] = con
	serr := dbSess.Save(c.Request(), c.Response().Writer)
	if serr != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("unable to save session %s", serr))
	}
	return c.JSON(http.StatusOK, dbSess.Values["db"])
}
