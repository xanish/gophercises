package phone

import (
	"database/sql"
	"fmt"
	"strings"
)

type Contacts struct {
	db *sql.DB
}

type Contact struct {
	Id     int
	Number string
}

func Open(driverName, dataSource string) (*Contacts, error) {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}

	return &Contacts{db}, nil
}

func (c *Contacts) Close() error {
	return c.db.Close()
}

func (c *Contacts) All() ([]Contact, error) {
	rows, err := c.db.Query("SELECT * FROM contacts ORDER BY id ASC")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch contacts: %v", err)
	}
	defer rows.Close()

	var ret []Contact
	for rows.Next() {
		var ct Contact
		if err := rows.Scan(&ct.Id, &ct.Number); err != nil {
			return nil, err
		}
		ret = append(ret, ct)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}

func (c *Contacts) Find(number string) ([]Contact, error) {
	rows, err := c.db.Query("SELECT * FROM contacts WHERE number = $1 ORDER BY id ASC", number)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch contacts matching %s: %v", number, err)
	}
	defer rows.Close()

	var ret []Contact
	for rows.Next() {
		var ct Contact
		if err := rows.Scan(&ct.Id, &ct.Number); err != nil {
			return nil, err
		}
		ret = append(ret, ct)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}

func (c *Contacts) Create(numbers []string) error {
	var query strings.Builder

	query.WriteString("INSERT INTO `contacts` (`number`) VALUES ")
	for _ = range numbers {
		query.WriteString("(?), ")
	}

	_, err := c.db.Exec(strings.Trim(query.String(), ","), numbers)
	if err != nil {
		return fmt.Errorf("failed to create contacts: %w", err)
	}

	return nil
}

func (c *Contacts) Update(contact Contact) error {
	_, err := c.db.Exec("UPDATE `contacts` SET number = $1 WHERE id = $2", contact.Number, contact.Id)
	if err != nil {
		return fmt.Errorf("failed to update contact %s: %w", contact.Number, err)
	}

	return nil
}

func (c *Contacts) Delete(numbers []string) error {
	_, err := c.db.Exec("DELETE FROM `contacts` WHERE number in ($1)", strings.Join(numbers, ","))
	if err != nil {
		return fmt.Errorf("failed to delete contacts: %w", err)
	}

	return nil
}
