package main
 
import (
	"bytes"
	"encoding/json"
	"fmt"
	r "github.com/dancannon/gorethink"
	"io/ioutil"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)
 
type Zip struct {
	Id    string    `json:"_id" gorethink:"id" bson:"_id"`
	City  string    `json:"city"`
	State string    `json:"state"`
	Pop   int       `json:"int"`
	Loc   []float64 `json:"loc"`
}
 
func Getenv(key, default_ string) string {
	if val := os.Getenv(key); val != "" {
		return val
	} else {
		return default_
	}
}
 
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
 
	rSession, err := r.Connect(map[string]interface{}{
		"address":     Getenv("RETHINKDB_URI", "localhost:28015"),
		"database":    "test",
		"maxIdle":     100,
		"idleTimeout": time.Second * 10,
	})
	if err != nil {
		log.Fatal(err)
	}
 
	r.Db("test").TableDrop("zips").RunWrite(rSession)
	r.Db("test").TableCreate("zips").RunWrite(rSession)
 
	mSession, err := mgo.Dial(Getenv("MONGODB_URI", "localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	mSession.DB("test").C("zips").DropCollection()
 
	f, err := os.Open("zips.json")
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
 
	// this is used to control concurrency
	c := make(chan bool, 100)
 
	rethinkWriteTest(rSession, json.NewDecoder(bytes.NewReader(data)), c)
	mongoWriteTest(mSession, json.NewDecoder(bytes.NewReader(data)), c)
	rethinkReadTest(rSession)
	mongoReadTest(mSession)
}
 
func rethinkWriteTest(session *r.Session, dec *json.Decoder, c chan bool) {
	var wg sync.WaitGroup
	start := time.Now()
	for {
		var zip Zip
		if err := dec.Decode(&zip); err != nil {
			break
		}
 
		wg.Add(1)
		c <- true
		go func(zip *Zip) {
			defer wg.Done()
			_, err := r.Table("zips").Insert(zip).RunWrite(session)
			if err != nil {
				log.Fatal(err)
			}
			<-c
		}(&zip)
	}
 
	wg.Wait()
	log.Printf("RethinkDB Write: %v\n", time.Since(start))
}
 
func mongoWriteTest(session *mgo.Session, dec *json.Decoder, c chan bool) {
	var wg sync.WaitGroup
	start := time.Now()
	for {
		var zip Zip
		if err := dec.Decode(&zip); err != nil {
			break
		}
 
		wg.Add(1)
		c <- true
		go func(zip *Zip) {
			defer wg.Done()
			session.DB("test").C("zips").Insert(zip)
			<-c
		}(&zip)
	}
 
	wg.Wait()
	log.Printf("MongoDB   Write: %v\n", time.Since(start))
}
 
func rethinkReadTest(session *r.Session) {
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < 10; i++ {
		lower := fmt.Sprintf("%d0000", i)
		upper := fmt.Sprintf("%d0000", i+1)
		wg.Add(1)
		go func(lower, upper string) {
			defer wg.Done()
			var zips []Zip
			rows, err := r.Table("zips").Between(lower, upper).Run(session)
			if err != nil {
				log.Fatal(err)
			}
			rows.ScanAll(&zips)
		}(lower, upper)
	}
 
	wg.Wait()
	log.Printf("RethinkDB Read:  %v\n", time.Since(start))
}
 
func mongoReadTest(session *mgo.Session) {
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < 10; i++ {
		lower := fmt.Sprintf("%d0000", i)
		upper := fmt.Sprintf("%d0000", i+1)
		wg.Add(1)
		go func(lower, upper string) {
			defer wg.Done()
			var zips []Zip
			query := session.DB("test").C("zips").Find(
				bson.M{"_id": bson.M{"$gt": lower, "$lt": upper}})
			if err := query.All(&zips); err != nil {
				log.Fatal(err)
			}
		}(lower, upper)
	}
 
	wg.Wait()
	log.Printf("MongoDB   Read:  %v\n", time.Since(start))
}