package main

import (
	"encoding/gob"
	"gotest/contacts"
	"gotest/db"
	"gotest/session"
	"gotest/user"

	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()
	gob.Register(db.Connection{})
	app.Use(session.UserSessions())
	app.POST("/connect", user.Connect)
	app.GET("/myinfo", user.GetInfo)
	app.GET("/contacts/:id", contacts.GetContact)
	app.POST("/contacts", contacts.CreateContact)
	app.DELETE("contacts/:id", contacts.DeleteContact)
	app.PUT("contacts/:id", contacts.UpdateContact)
	app.Start(":8080")
}
