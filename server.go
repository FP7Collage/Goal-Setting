package main

import (
	"flag"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/encoder"
	"log"
	"net/http"
)

func main() {
	var host string
	var port int
	flag.StringVar(&host, "host", "", "The address to bind to")
	flag.IntVar(&port, "port", 9999, "The port to listen at")

	m := martini.Classic()
	m.Use(func(c martini.Context, w http.ResponseWriter) {
		c.MapTo(encoder.JsonEncoder{PrettyPrint: false}, (*encoder.Encoder)(nil))
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	})

	m.Get("/users", GetUsers)
	m.Post("/users", binding.Bind(User{}), CreateUser)
	m.Get("/users/:id", GetUser)
	m.Delete("/users/:id", DeleteUser)
	m.Get("/users/:id/teams", GetUserTeams)
	m.Get("/users/:id/goals", GetUserGoals)

	m.Get("/teams", GetTeams)
	m.Post("/teams", binding.Bind(Team{}), CreateTeam)
	m.Get("/teams/:id", GetTeam)
	m.Delete("/teams/:id", DeleteTeam)

	m.Put("/users/:uid/teams/:tid", AddUserToTeam)
	m.Delete("/users/:uid/teams/:tid", DeleteUserFromTeam)

	m.Post("/goals", binding.Bind(Goal{}), CreateGoal)
	m.Get("/goals/:id", GetGoal)
	m.Post("/goals/:id", binding.Bind(GoalCompleteParams{}), CompleteGoal)

	log.Fatal(http.ListenAndServe(":8080", m))
	m.Run()
}
