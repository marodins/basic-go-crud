package contacts

import (
	"gotest/db"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

func GetContact(c echo.Context) error {
	query := `SELECT * FROM contacts WHERE id = $1;`
	sess := c.Get("database").(*sessions.Session)
	con := sess.Values["db"]
	if con == nil {
		return c.JSON(http.StatusForbidden, "no connection")
	}
	id := c.Param("id")
	db, err0 := con.(db.Connection).GetConnection()

	if err0 != nil {
		return c.JSON(http.StatusForbidden, "can't establish connection")
	}
	contact := Contact{}
	if err := db.QueryRow(query, id).Scan(&contact.First, &contact.Last, &contact.Phone, &contact.id); err != nil {
		return c.JSON(400, map[string]string{"msg": "unable to initialize contact", "error": err.Error()})
	}

	return c.JSON(200, map[string]string{
		"first": contact.First,
		"last":  contact.Last,
		"phone": contact.Phone,
	})

}

func CreateContact(c echo.Context) error {
	q := `CREATE TABLE IF NOT EXISTS contacts(
			first varchar(50) NOT NULL,
			last varchar(50) NOT NULL,
			phone varchar(50) NOT NULL,
			id serial primary key
		);`

	contactq := `INSERT INTO contacts(first, last, phone) VALUES ($1, $2, $3) RETURNING id;`

	sess := c.Get("database").(*sessions.Session)
	con := sess.Values["db"]
	var contact Contact
	if con == nil {
		return c.JSON(http.StatusBadRequest, "no connection")
	}

	db, err0 := con.(db.Connection).GetConnection()
	if err0 != nil {
		return c.JSON(http.StatusForbidden, "can't establish connection")
	}
	if c.Bind(&contact) != nil {
		return c.JSON(http.StatusBadRequest, "incorrect data submitted")
	}
	_, err1 := db.Exec(q)

	if err1 != nil {
		return c.JSON(401, map[string]string{"error": err1.Error(), "message": "unable to create table"})
	}
	res := db.QueryRow(contactq, contact.First, contact.Last, contact.Phone)

	var id int
	res.Scan(&id)
	return c.JSON(200, map[string]int{"id": id})
}

func DeleteContact(c echo.Context) error {
	q := `DELETE from contacts WHERE id = $1;`
	sess := c.Get("database").(*sessions.Session)
	con := sess.Values["db"]
	if con == nil {
		return c.JSON(http.StatusForbidden, "no connection")
	}
	id := c.Param("id")
	db, err0 := con.(db.Connection).GetConnection()
	if err0 != nil {
		return c.JSON(http.StatusForbidden, "can't establish connection")
	}

	if _, err := db.Exec(q, id); err != nil {
		return c.JSON(400, "cannot delete that contact")
	}
	return c.NoContent(204)
}

func UpdateContact(c echo.Context) error {
	q := `UPDATE contacts
	SET first = $1, last = $2, phone = $3
	WHERE id = $4;`
	sess := c.Get("database").(*sessions.Session)
	con := sess.Values["db"]
	if con == nil {
		return c.JSON(http.StatusForbidden, "no connection")
	}
	id := c.Param("id")
	db, err0 := con.(db.Connection).GetConnection()
	contact := Contact{}
	c.Bind(&contact)
	if err0 != nil {
		return c.JSON(http.StatusForbidden, "can't establish connection")
	}

	if _, err := db.Exec(q, contact.First, contact.Last, contact.Phone, id); err != nil {
		return c.JSON(400, "cannot delete that contact")
	}
	return c.NoContent(204)
}
