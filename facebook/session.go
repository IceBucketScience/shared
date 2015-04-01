package facebook

import (
	fb "github.com/huandu/facebook"
)

type Session struct {
	fbSession *fb.Session
}

func CreateSession(accessToken string) *Session {
	return &Session{fbSession: fbClient.Session(accessToken)}
}

func (session *Session) GetInfo() (*Person, error) {
	res, fbErr := session.fbSession.Get("/me", nil)
	if fbErr != nil {
		return nil, fbErr
	}

	var info Person
	decodeErr := res.Decode(&info)
	if decodeErr != nil {
		return nil, decodeErr
	}

	return &info, nil
}

type Person struct {
	UserId string `facebook:"id"`
	Name   string
}

func (session *Session) GetFriends() ([]*Person, error) {
	res, fbErr := session.fbSession.Get("/me/friends", nil)
	if fbErr != nil {
		return nil, fbErr
	}

	var friendsList []*Person
	decodeErr := res.DecodeField("data", &friendsList)
	if decodeErr != nil {
		return nil, decodeErr
	}

	return friendsList, nil
}
