package graph

import (
	"github.com/jmcvetta/neoism"
)

type Person struct {
	Node *neoism.Node
}

func CreatePerson(userId string, name string) (*Person, error) {
	node, err := db.CreateNode(neoism.Props{"fbId": userId, "name": name})
	if err != nil {
		return nil, err
	}

	//TODO: check for errors on adding labels

	node.AddLabel("Person")

	return &Person{Node: node}, nil
}

func GetPerson(userId string) (*Person, error) {
	res := []struct {
		Person neoism.Node
	}{}

	err := db.Cypher(&neoism.CypherQuery{
		Statement:  `MATCH (p:Person) WHERE p.fbId = {userId} RETURN p`,
		Parameters: neoism.Props{"userId": userId},
		Result:     &res,
	})
	if err != nil {
		return nil, err
	} else if len(res) < 1 {
		return nil, nil
	}

	return &Person{Node: &res[0].Person}, nil
}

type person interface {
}
