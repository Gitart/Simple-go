package main
 
import (
	r "github.com/dancannon/gorethink"
	"log"
)
 
type upd struct {
	Zzz []string `gorethink:"zzz"`
}
 
func main() {
	session, err := r.Connect(map[string]interface{}{
		"address":  "rethinkdb:28015",
		"database": "test",
	})
	if err != nil {
		log.Fatalln(err)
	}
 
	res1, e1 := r.Table("tv_shows").Insert(map[string]interface{}{
		"name": "aaa",
		"zzz":  []string{"ccc", "ddd"},
	}).RunWrite(session)
	if e1 != nil {
		log.Fatalln(e1)
	}
 
	term := r.Table("tv_shows").Get(res1.GeneratedKeys[0])
	q := map[string]r.RqlTerm{"zzz": r.Row.Field("zzz").Append("ppp1")}
	_, e2 := term.Update(q).RunWrite(session)
	if e2 != nil {
		log.Fatalln(e2)
	}
 
	_, e3 := r.Table("tv_shows").Get(res1.GeneratedKeys[0]).Delete().RunWrite(session)
	if e3 != nil {
		log.Fatalln(e3)
	}
}

