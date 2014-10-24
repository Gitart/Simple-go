package main

import (
	"fmt"
	"log"

	r "github.com/dancannon/gorethink"
)

/*
type upd struct {
	Zzz []string `gorethink:"zzz"`
}





*/

// Тип
type Person struct {
	Id        string `gorethink:"id, omitempty"`
	FirstName string `gorethink:"first_name"`
	LastName  string `gorethink:"last_name"`
	Gender    string `gorethink:"gender"`
}

// Функция раоботы с базой
func main() {

	session, err := r.Connect(r.ConnectOpts{Address: "localhost:28015", Database: "test"})

	if err != nil {
		log.Fatalln(err)
	}

	// Вставка операций для примера
	//r.Table("tv_shows").Insert(map[string]interface{}{"name": "aaa", "zzz": []string{"ccc", "ddd"}}).RunWrite(session)
	//r.Table("tv_shows").Insert(map[string]interface{}{"name": "Пример"}).RunWrite(session)
	//r.table("tv_shows").Insert({"name": "Пример yjdjcn"}).RunWrite(session)
	//r.table("posts").insert(map[string]interface{}{ id: 1, title: "Lorem ipsum", content: "Dolor sit amet"}).RunWrite(session)

	//
	/*
		r.Db("test").Table("post").Insert(Person{"1", "John", "Smith", "M"}).Run(session)
		fmt.Println("Вставка готова")

		r.Db("test").Table("Test").Insert(Person{"1", "John", "Smith", "M"}).Run(session)
		r.Db("test").Table("Test").Insert(Person{"2", "John2", "Smith2", "M2"}).Run(session)

	*/

	//r.Db("test").Table("Test").Insert({"sddd": "2223333ee"}).Run(session)

	r.Db("test").TableDrop("table").Run(session)

	response, err := r.Db("test").TableCreate("table").RunWrite(session)

	if err != nil {
		log.Fatalf("Error creating table: %s", err)
	}

	//fmt.Printf("%d table created", response.Created)

	//fmt.Println("Вставка готова тест")

	//Dins("Test", "id4:2333,       Name:'ssdddffff-00'")
	//Dins("Test", "id3:A2333423,  Name:'Номер входного парметра'")

	fmt.Println("Вставка функций")

	res, err := r.Expr("Hello World").Run(session)
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = res.One(&response)
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println(response)

}

func Dins(Tabname string, CommInsert string) {
	session, err := r.Connect(r.ConnectOpts{Address: "10.0.20.5:28015", Database: "test"})

	if err != nil {
		//log.Fatalln(err.Error())
		fmt.Println(err.Error())
	}

	r.Table(Tabname).Insert(CommInsert).Run(session)

}
