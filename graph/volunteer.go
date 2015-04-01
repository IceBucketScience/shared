package graph

import (
	"github.com/jmcvetta/neoism"
)

type Volunteer struct {
	Node *neoism.Node
}

func CreateVolunteer(userId string, name string, accessToken string) (*Volunteer, error) {
	person, personErr := CreatePerson(userId, name)
	if personErr != nil {
		return nil, personErr
	}

	node := person.Node

	//TODO: check for errors on adding labels and properties

	//TODO: due to a bug in neoism, right now I have to reset all properties. Once the bug is fixed,
	//then I can use the following code to set the extra Volunteer properties:
	/*person.SetProperty("accessToken", accessToken)
	person.SetProperty("isIndexed", false)*/

	//Temporary code to set properties until bug is fixed
	props, _ := node.Properties()

	props["accessToken"] = accessToken
	props["isIndexed"] = false
	node.SetProperties(props)
	//End temporary code

	node.AddLabel("Volunteer")

	return &Volunteer{Node: node}, nil
}

func (volunteer *Volunteer) MarkAsIndexed() {
	node := volunteer.Node
	//TODO: check for errors on adding labels and properties
	//TODO: once bug is fixed, use volunteer.SetProperty
	props, _ := node.Properties()

	props["isIndexed"] = true
	node.SetProperties(props)
}

func (volunteer *Volunteer) hasBeenLinked(userId string) {

}

func FindIndexedVolunteer(userId string) (*neoism.Node, error) {
	res := []struct {
		Volunteer neoism.Node
	}{}

	err := db.Cypher(&neoism.CypherQuery{
		Statement:  `MATCH (v:Volunteer) WHERE v.fbId = {userId} AND v.isIndexed = true RETURN v`,
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
