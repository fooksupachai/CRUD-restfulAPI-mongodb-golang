package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Developer struct
type Developer struct {
	ID   string `bson:"id"`
	Name string `bson:"name"`
}

func getDevelopers() []Developer {
	session, err := mgo.Dial("mongodb://127.0.0.1:27017/testDB")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("testDB").C("developers")
	developers := []Developer{}
	err = c.Find(nil).All(&developers)
	if err != nil {
		log.Fatal(err)
	}
	return developers
}

func getDeveloper(id string) Developer {
	session, err := mgo.Dial("mongodb://127.0.0.1:27017/testDB")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("testDB").C("developers")
	developer := Developer{}
	err = c.Find(bson.M{"id": id}).One(&developer)
	if err != nil {
		log.Fatal(err)
	}
	return developer

}

func removeDeveloper(id string) Developer {
	session, err := mgo.Dial("mongodb://127.0.0.1:27017/testDB")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("testDB").C("developers")
	developer := Developer{}
	err = c.Remove(bson.M{"id": id})
	if err != nil {
		log.Fatal(err)
	}
	return developer

}

func createDeveloper(e Developer) error {
	session, err := mgo.Dial("mongodb://127.0.0.1:27017/testDB")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("testDB").C("developers")
	return c.Insert(e)
}

func updateDeveloper(e Developer) error {
	session, err := mgo.Dial("mongodb://127.0.0.1:27017/testDB")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	newID := "0000"
	selector := bson.M{"id": e.ID}
	updator := bson.M{"$set": bson.M{"id": newID}}

	c := session.DB("testDB").C("developers")
	return c.Update(selector, updator)
}

func getAllDevelopers(res http.ResponseWriter, req *http.Request) {
	developers := getDevelopers()
	json.NewEncoder(res).Encode(developers)
}

func getOneDeveloper(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	developer := getDeveloper(params["id"])
	json.NewEncoder(res).Encode(developer)
}

func postDeveloper(res http.ResponseWriter, req *http.Request) {
	var item Developer
	json.NewDecoder(req.Body).Decode(&item)
	createDeveloper(item)
}

func putDeveloper(res http.ResponseWriter, req *http.Request) {
	var item Developer
	json.NewDecoder(req.Body).Decode(&item)
	updateDeveloper(item)
}

func DeleteDeveloper(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	removeDeveloper(params["id"])
	res.Write([]byte("OK"))
}

func handleRequest() {
	router := mux.NewRouter()
	router.HandleFunc("/getAllDevelopers", getAllDevelopers).Methods("GET")
	router.HandleFunc("/getDeveloper/{id}", getOneDeveloper).Methods("GET")
	router.HandleFunc("/postDeveloper", postDeveloper).Methods("POST")
	router.HandleFunc("/putDeveloper", putDeveloper).Methods("PUT")
	router.HandleFunc("/removeDeveloper/{id}", DeleteDeveloper).Methods("DELETE")
	http.ListenAndServe(":8080", router)
}

func main() {
	handleRequest()
}
