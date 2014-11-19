package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/encoder"
	"github.com/nu7hatch/gouuid"
)

type Goal struct {
	Id                  string `json:"id"`
	Type                string `json:"type" binding:"required"`
	Target              string `json:"target" binding:"required"`
	TargetType          string `json:"target_type"`
	Start               int    `json:"start"`
	Length              int    `json:"length" binding:"required"`
	Name                string `json:"name" binding:"required"`
	Content             string `json:"content" binding:"required"`
	Keywords            string `json:"keywords" binding:"required"`
	ActiveAfter         string `json:"active_after"`
	Reward              string `json:"reward" binding:"required"`
	Difficulty          int    `json:"difficulty" binding:"required"`
	NumberOfCompletions int    `json:"number_of_completions" binding:"required"`
	State               string `json:"state"`
}

type GoalCompleteParams struct {
	UserId string `json:"user_id" binding:"required"`
}

var goals = make(map[string]*Goal)

func CreateGoal(enc encoder.Encoder, goal Goal) (int, []byte) {
	goalId, err := uuid.NewV4()
	if err != nil {
		return 500, encoder.Must(enc.Encode(""))
	}
	goal.Id = goalId.String()
	if goal.TargetType == "user" {
		user, ok := users[goal.Target]
		if !ok {
			return 400, encoder.Must(enc.Encode("Target does not exist"))
		}
		user.Goals[goal.Id] = &goal
	} else if goal.TargetType == "team" {
		team, ok := teams[goal.Target]
		if !ok {
			return 400, encoder.Must(enc.Encode("Target does not exist"))
		}
		team.Goals[goal.Id] = &goal
	} else {
		return 400, encoder.Must(enc.Encode("Unknown target type"))
	}

	goals[goal.Id] = &goal
	return 201, encoder.Must(enc.Encode(goal.Id))
}

func GetGoal(enc encoder.Encoder, params martini.Params) (int, []byte) {
	goal, ok := goals[params["id"]]
	if ok {
		return 200, encoder.Must(enc.Encode(goal))
	}
	return 404, encoder.Must(enc.Encode())
}

func CompleteGoal(enc encoder.Encoder, params martini.Params, jsonParams GoalCompleteParams) (int, []byte) {
	goal, ok := goals[params["id"]]
	if !ok {
		return 404, encoder.Must(enc.Encode())
	}
	if _, ok := users[jsonParams.UserId]; !ok {
		return 400, encoder.Must(enc.Encode("User does not exist"))
	}
	goal.NumberOfCompletions--
	if goal.NumberOfCompletions > 0 {
		goal.State = "progress"
		return 200, encoder.Must(enc.Encode(goal.State)) //what should be returned here?
	}

	goal.State = "completed"
	return 200, encoder.Must(enc.Encode(goal.Reward))
}
