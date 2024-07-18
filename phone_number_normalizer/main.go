package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/rodaine/table"
	"github.com/xanish/gophercises/phone_number_normalizer/normalize"
	"github.com/xanish/gophercises/phone_number_normalizer/phone"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "dev"
	password = "secret"
	dbname   = "gophercises_phone"
)

var seeds = []string{
	"1234567890",
	"123 456 7891",
	"(123) 456 7892",
	"(123) 456-7893",
	"123-456-7894",
	"123-456-7890",
	"1234567892",
	"(123)456-7892",
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	ph, err := phone.NewPhone("postgres", psqlInfo)
	must(err)

	err = ph.Contacts.Reset()
	must(err)

	err = ph.Contacts.Init()
	must(err)

	err = ph.Contacts.Create(seeds)
	must(err)

	contacts, err := ph.Contacts.All()
	must(err)
	fmt.Println("Contacts before normalization")
	printContacts(contacts)

	for _, c := range contacts {
		number := normalize.Normalize(c.Number)
		if number == c.Number {
			continue
		}

		existing, err := ph.Contacts.Find(number)
		must(err)
		if existing != nil {
			must(ph.Contacts.Delete([]string{c.Number}))
		} else {
			c.Number = number
			must(ph.Contacts.Update(c))
		}
	}

	contacts, err = ph.Contacts.All()
	must(err)
	fmt.Println("\nContacts after normalization")
	printContacts(contacts)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func printContacts(contacts []phone.Contact) {
	tbl := table.New("#", "Number")
	for _, c := range contacts {
		tbl.AddRow(c.Id, c.Number)
	}
	tbl.Print()
}
