package router

import (
	db "CBCTF/internel/repo"
	"CBCTF/internel/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HomePage(ctx *gin.Context) {
	DB := db.DB.WithContext(ctx)

	data := gin.H{
		"upcoming":   []gin.H{},
		"stats":      []gin.H{},
		"scoreboard": []gin.H{},
	}
	contests, count, ok, _ := db.InitContestRepo(DB).GetAll(-1, -1, true, 0, false)
	if ok {
		for i := 0; i < func() int {
			if len(contests) > 3 {
				return 3
			}
			return len(contests)
		}(); i++ {
			contest := contests[i]
			info := gin.H{
				"name":     contest.Name,
				"start":    contest.Start,
				"duration": contest.Duration.Seconds(),
				"users":    len(contest.Users),
				"teams":    len(contest.Teams),
				"avatar":   contest.Avatar,
			}
			data["upcoming"] = append(data["upcoming"].([]gin.H), info)
		}
	}
	data["stats"] = append(data["stats"].([]gin.H), gin.H{"label": "CTF Events", "value": count})
	count, _, _ = db.InitUserRepo(DB).Count(true, true)
	data["stats"] = append(data["stats"].([]gin.H), gin.H{"label": "Activate CTFers", "value": count})
	count, _, _ = db.InitChallengeRepo(DB).Count()
	data["stats"] = append(data["stats"].([]gin.H), gin.H{"label": "Challenges", "value": count})
	count, _, _ = db.InitSubmissionRepo(DB).Count()
	data["stats"] = append(data["stats"].([]gin.H), gin.H{"label": "Submissions", "value": count})
	users, _, _, _ := service.GetUserRanking(DB, 5, 0)
	for _, user := range users {
		data["scoreboard"] = append(data["scoreboard"].([]gin.H), gin.H{
			"name":    user.Name,
			"score":   user.Score,
			"solved":  user.Solved,
			"country": user.Country,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": data})
}
