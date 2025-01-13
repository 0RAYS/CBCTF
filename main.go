package main

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/router"
	"fmt"
)

func init() {
	i18n.Init()
	config.Init()
	log.Init()
	db.Init()
}

func main() {
	ip, port := config.Env.GetString("gin.ip"), config.Env.GetString("gin.port")
	if err := router.InitRouters().Run(fmt.Sprintf("%s:%s", ip, port)); err != nil {
		log.Logger.Panicf("Failed to start: %s", err)
	}
	//db.CreateUser("user1", "1", "1@test.com")
	//db.CreateContest("contest1")
	//db.CreateTeam(1, "team1", 1)
	//
	//db.UpdateUser(1, map[string]interface{}{"name": "user1_updated", "id": 2})
	//db.UpdateContest(1, map[string]interface{}{"name": "contest1_updated", "id": 1})
	//db.UpdateTeam(1, map[string]interface{}{"name": "team1_updated", "id": 1})
	//
	//user1, ok, msg := db.GetUserByID(1, true)
	//fmt.Println(user1.Teams, msg, ok)
	//team1, ok, msg := db.GetTeamByID(1, true)
	//fmt.Println(team1.Contests, team1.Users, msg, ok)
	//contest1, ok, msg := db.GetContestByID(1, true)
	//fmt.Println(contest1.Teams, msg, ok)
	//
	//db.CreateUser("user2", "2", "2@test.com")
	//db.CreateContest("contest2")
	//db.CreateTeam(2, "team2", 2)
	//res, msg := db.JoinTeam(2, 1, 1)
	//fmt.Println(res, msg)
	//res, msg := db.LeaveTeam(2, 1, 1)
	//fmt.Println(res, msg)
	//db.DeleteContest(2)
	//db.DeleteTeam(2)
	//db.DeleteUser(2)
	//res := db.ClearEmptyTeam()

	//user, _, _ := db.GetUserByID(1)
	//fmt.Println(user)
	//data := utils.TidyRetData(user)
	//fmt.Println(data)

	//users, _, _ := db.GetUsers(0, 0)
	//fmt.Println(map[string]interface{}{"users": "test"})
	//data := utils.TidyRetData(users)
	//fmt.Println(data, len(data))
	//fmt.Println(data[0])
	//user, _, _ := db.GetUserByID(1)
	//utils.TidyRetData(user)
}
