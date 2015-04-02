package graph

import (
	"github.com/jmcvetta/neoism"
)

type Person struct {
	node *neoism.Node

	FbId, Name string
}

func CreatePerson(userId string, name string) (*Person, error) {
	node, err := db.CreateNode(neoism.Props{"fbId": userId, "name": name})
	if err != nil {
		return nil, err
	}

	//TODO: check for errors on adding labels

	node.AddLabel("Person")

	return &Person{node: node, FbId: userId, Name: name}, nil
}

func (person *Person) GetDb() *neoism.Database {
	return person.node.Db
}

func (person *Person) AddFriendshipWith(friend *Person) error {
	/*_, err := person.node.Relate("FRIENDS", friend.node.Id(), neoism.Props{})
	return err*/

	res := []struct {
		F neoism.Relationship
	}{}

	err := db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (p1:Person), (p2:Person)
            WHERE p1.fbId = {p1Id} AND p2.fbId = {p2Id}
            MERGE (p1)-[f:FRIENDS]-(p2)
            RETURN f
        `,
		Parameters: neoism.Props{"p1Id": person.FbId, "p2Id": friend.FbId},
		Result:     &res,
	})
	if err != nil {
		return err
	}

	return nil
}

func getPersonFromNode(node *neoism.Node) (*Person, error) {
	props, err := node.Properties()
	if err != nil {
		return nil, err
	}

	return &Person{node: node, FbId: props["fbId"].(string), Name: props["name"].(string)}, nil
}

func GetPerson(userId string) (*Person, error) {
	res := []struct {
		Person neoism.Node `json:"p"`
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

	// adds a db object to each node
	for index, personData := range res {
		personData.Person.Db = db
		res[index] = personData
	}

	return getPersonFromNode(&res[0].Person)
}
