package util

/**
*    WORK IN PROGRESS
*/

import (
	"encoding/json"
	"fmt"
)

func SyncUsers(){

	users := PullUsersFromAD()

	userList := ParseUsers(users)

	GenerateUserAcessPolicies(userList)

	// Insert policies into database

}

func PullUsersFromAD() string {

	// Logic to pull from Graph API

	users := ""

	fmt.Println(users)

	return ""
}

func ParseUsers(users string) map[string]string {

	var userList map[string]string
	json.Unmarshal([]byte(users), &userList)

	return userList
}

func GenerateUserAcessPolicies(userList map[string]string){

	// Generates Casbin access control policies

}