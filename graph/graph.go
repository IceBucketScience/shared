package graph

import (
	"log"

	"github.com/jmcvetta/neoism"
)

var db *neoism.Database

func InitGraph(dbUrl string) {
	dbObj, err := neoism.Connect(dbUrl)
	if err != nil {
		log.Fatalln(err)
	}

	db = dbObj
}
