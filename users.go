package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/encoder"
)

type User struct {
	Id    string `json:"id" binding:"required"`
	Teams map[string]string
	Goals map[string]*Goal
}

type Team struct {
	Id    string `json:"id" binding:"required"`
	Users map[string]string
	Goals map[string]*Goal
}

var users = make(map[string]*User)
var teams = make(map[string]*Team)

func CreateUser(enc encoder.Encoder, user User) (int, []byte) {
	user.Teams = make(map[string]string)
	user.Goals = make(map[string]*Goal)
	users[user.Id] = &user
	return 201, encoder.Must(enc.Encode(user.Id))
}

func CreateTeam(enc encoder.Encoder, team Team) (int, []byte) {
	team.Users = make(map[string]string)
	team.Goals = make(map[string]*Goal)
	teams[team.Id] = &team
	return 201, encoder.Must(enc.Encode(team.Id))
}

func GetUsers(enc encoder.Encoder) (int, []byte) {
	return 200, encoder.Must(enc.Encode(users))
}

func GetTeams(enc encoder.Encoder) (int, []byte) {
	return 200, encoder.Must(enc.Encode(teams))
}

func GetUser(enc encoder.Encoder, params martini.Params) (int, []byte) {
	if user, ok := users[params["id"]]; ok {
		return 200, encoder.Must(enc.Encode(user))
	}
	return 404, encoder.Must(enc.Encode())
}

func GetTeam(enc encoder.Encoder, params martini.Params) (int, []byte) {
	if team, ok := teams[params["id"]]; ok {
		return 200, encoder.Must(enc.Encode(team))
	}
	return 404, encoder.Must(enc.Encode())
}

func DeleteUser(params martini.Params) (int, string) {
	if user, ok := users[params["id"]]; ok {
		for key := range user.Teams {
			delete(teams[key].Users, user.Id)
		}
		delete(users, user.Id)
		return 204, ""
	}
	return 404, ""
}

func DeleteTeam(params martini.Params) (int, string) {
	if team, ok := teams[params["id"]]; ok {
		for key := range team.Users {
			delete(users[key].Teams, team.Id)
		}
		delete(teams, team.Id)
		return 204, ""
	}
	return 404, ""
}

func GetUserTeams(enc encoder.Encoder, params martini.Params) (int, []byte) {
	if user, ok := users[params["id"]]; ok {
		return 200, encoder.Must(enc.Encode(user.Teams))
	}
	return 404, encoder.Must(enc.Encode())
}

func GetUserGoals(enc encoder.Encoder, params martini.Params) (int, []byte) {
	if user, ok := users[params["id"]]; ok {
		return 200, encoder.Must(enc.Encode(user.Goals))
	}
	return 404, encoder.Must(enc.Encode())
}

func AddUserToTeam(params martini.Params) (int, string) {
	user, user_ok := users[params["uid"]]
	team, team_ok := teams[params["tid"]]
	if !user_ok || !team_ok {
		return 404, ""
	}
	user.Teams[team.Id] = team.Id
	team.Users[user.Id] = user.Id
	return 204, ""
}

func DeleteUserFromTeam(params martini.Params) (int, string) {
	user, user_ok := users[params["uid"]]
	team, team_ok := teams[params["tid"]]
	if !user_ok || !team_ok {
		return 404, ""
	}
	delete(user.Teams, team.Id)
	delete(team.Users, user.Id)
	return 204, ""
}
