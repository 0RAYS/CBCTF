package router

import (
	"CBCTF/internel/middleware"
	"CBCTF/internel/resp"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetContest(ctx *gin.Context) {
	contest := middleware.GetContest(ctx)
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": resp.GetContestResp(contest)})
}
