package main
 
import (
	r "github.com/dancannon/gorethink"
	"log"
)
 
type valT struct {
	V float64
}
 
func main() {
	dbs, err := r.Connect(map[string]interface{}{
		"address":  "localhost:28015",
		"database": "test",
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
 
	rows, err := r.Db("test").Table("values").Run(dbs)
	if err != nil {
		log.Fatalln(err.Error())
	}
	var total float64 = 0.0
	for rows.Next() {
		var row valT
		if err := rows.Scan(&row); err != nil {
			log.Fatalln(err)
		}
		total += row.V
	}
	println(total)
}