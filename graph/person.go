package graph

import (
	"errors"

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

func (person *Person) addRelationshipTo(fbId string, relName string) error {
	res := []struct {
		F neoism.Relationship
	}{}

	err := db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (p1:Person), (p2:Person)
            WHERE p1.fbId = {p1Id} AND p2.fbId = {p2Id}
            MERGE (p1)-[f:` + relName + `]-(p2)
            RETURN f
        `,
		Parameters: neoism.Props{"p1Id": person.FbId, "p2Id": fbId},
		Result:     &res,
	})
	if err != nil {
		return err
	} else if len(res) < 1 {
		return errors.New("no new relationship created or existing relationship found")
	}

	return nil
}

func (person *Person) AddFriendshipWith(friendId string) error {
	return person.addRelationshipTo(friendId, "FRIENDS")
}

func (person *Person) MarkAsLinkedTo(volunteer *Volunteer) error {
	return person.addRelationshipTo(volunteer.FbId, "LINKED")
}

func (person *Person) IsLinkedTo(volunteer *Volunteer) (bool, error) {
	res := []struct {
		L neoism.Relationship
	}{}

	err := db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (p1:Person)-[l:LINKED]-(p2:Person)
            WHERE p1.fbId = {p1Id} AND p2.fbId = {p2Id}
            RETURN l
        `,
		Parameters: neoism.Props{"p1Id": person.FbId, "p2Id": volunteer.FbId},
		Result:     &res,
	})
	if err != nil {
		return false, err
	} else if len(res) < 1 {
		return false, nil
	}

	return true, nil
}

//TODO: switch return val to pointer
func (person *Person) GetFriends() (Graph, error) {
	res := []struct {
		P neoism.Node `json:"friends"`
	}{}

	err := db.Cypher(&neoism.CypherQuery{
		Statement:  `MATCH (p:Person)-[:FRIENDS]-(friends) WHERE p.fbId = {fbId} RETURN friends`,
		Parameters: neoism.Props{"fbId": person.FbId},
		Result:     &res,
	})
	if err != nil {
		return nil, err
	}

	friends := Graph{}

	for _, personData := range res {
		personNode := &personData.P
		personNode.Db = db

		person, getPersonErr := getPersonFromNode(personNode)
		if getPersonErr != nil {
			return nil, getPersonErr
		}

		friends[person.FbId] = person
	}

	return friends, nil
}

type Friendship struct {
	rel      *neoism.Relationship
	SourceId string
	TargetId string
}

type FriendshipRes struct {
	Id int `json:"id"`
}

func (friendship *Friendship) GetRelationshipId() int {
	return friendship.rel.Id()
}

func GetFriendshipsInNetwork(personId string) ([]*Friendship, error) {
	res := []struct {
		SourceId string              `json:"sourceId"`
		TargetId string              `json:"targetId"`
		F        neoism.Relationship `json:"f"`
	}{}

	err := db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (v:Person {fbId: {fbId}})-[f:FRIENDS]->(friend:Person) 
                RETURN v.fbId AS sourceId, f, friend.fbId AS targetId
            UNION MATCH (v:Person {fbId: {fbId}})-[:FRIENDS]-(p:Person)-[f:FRIENDS]->(friend:Person) 
                RETURN p.fbId AS sourceId, f, friend.fbId AS targetId
        `,
		Parameters: neoism.Props{"fbId": personId},
		Result:     &res,
	})
	if err != nil {
		return nil, err
	}

	friends := []*Friendship{}

	for _, friendshipData := range res {
		friendshipNode := &friendshipData.F
		friendshipNode.Db = db

		friendship := &Friendship{SourceId: friendshipData.SourceId, TargetId: friendshipData.TargetId, rel: friendshipNode}

		friends = append(friends, friendship)
	}

	return friends, nil
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
		P neoism.Node
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
		personData.P.Db = db
		res[index] = personData
	}

	return getPersonFromNode(&res[0].P)
}
