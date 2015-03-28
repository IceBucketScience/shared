package graph

import (
	"github.com/jmcvetta/neoism"
)

func CreateVolunteer(userId string, accessToken string) (*neoism.Node, error) {
	node, err := db.CreateNode(neoism.Props{"fbId": userId, "accessToken": accessToken})
	if err != nil {
		return nil, err
	}

	node.AddLabel("Person")
	node.AddLabel("Volunteer")

	return node, nil
}

func FindVolunteer(userId string) (*neoism.Node, error) {
	res := []struct {
		Volunteer neoism.Node
	}{}

	err := db.Cypher(&neoism.CypherQuery{
		Statement:  `MATCH (v:Volunteer) WHERE v.fbId = {userId} RETURN v`,
		Parameters: neoism.Props{"userId": userId},
		Result:     &res,
	})
	if err != nil {
		return nil, err
	} else if len(res) < 1 {
		return nil, nil
	}

	return &res[0].Volunteer, nil
}
