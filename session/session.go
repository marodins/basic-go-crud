package session

import (
	"github.com/labstack/echo/v4"
)

func UserSessions() echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, _ := Store.Get(c.Request(), "database")
			c.Set("database", sess)
			return next(c)
		}
	}

}
