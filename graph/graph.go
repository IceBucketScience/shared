package graph

import (
	"log"

	"github.com/jmcvetta/neoism"
)

var db *neoism.Database

func InitGraph(dbUrl string) {
	dbObj, err := neoism.Connect(dbUrl)
	if err != nil {
		//TODO: switch to returning err
		log.Fatalln(err)
	}

	db = dbObj
}
