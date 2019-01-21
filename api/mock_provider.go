package api

import (
	"fmt"
	"time"
)

type MockAuthProvider struct {
}

func NewMockAuthProvider() AuthProvider {
	return &MockAuthProvider{}
}

func (MockAuthProvider) Friends() ([]string, error) {
	result:= make ([]string,0)
	return result, nil
}

func (MockAuthProvider) FriendsData() ([]User, error) {
	result:= make ([]User, 0)
	return result, nil
}

func (MockAuthProvider) Auth() (User, error) {
	current_time:= fmt.Sprintf("",time.Now())
	user:= User{Uid:current_time, FirstName:"Mock user",LastName:current_time}
	return  user, nil
}

