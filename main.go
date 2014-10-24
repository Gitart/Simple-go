package main

import (
	"fmt"
	"log"

	r "github.com/dancannon/gorethink"
)

// Функция раоботы с базой
func main() {
	session, err := r.Connect(r.ConnectOpts{Address: "10.0.20.5:28015", Database: "test"})

	if err != nil {
		log.Fatalln(err)
	}

	//r.Db("test").Table("Test").Insert("{sddd: 2223333}").Run(session)

	response, err: = r.Db("test").Table("Test").Insert([{"idddw": "1235779"}]).RunWrite(session)

	if err != nil {
		log.Fatalf("Error creating index: %s", err)
	}

	fmt.Printf("%d index created", response.Created)
}
