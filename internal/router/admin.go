package router

import (
	"RayWar/internal/db"
	"RayWar/internal/log"
	"RayWar/internal/middleware"
	"RayWar/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

func CreateUser(ctx *gin.Context) {
	var createUserForm CreateUserForm
	trace := middleware.GetTraceID(ctx)
	if err := ctx.ShouldBind(&createUserForm); err == nil {
		username, password, email := createUserForm.Name, createUserForm.Password, createUserForm.Email
		user, ok, msg := db.CreateUser(username, password, email)
		if !ok {
			log.Logger.Infof("| %s | %s", trace, msg)
			ctx.JSONP(http.StatusInternalServerError, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
			return
		}
		log.Logger.Infof("| %s | %s:%d register", trace, user.Name, user.ID)
		ctx.JSONP(http.StatusOK, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": utils.TidyRetData(user, "password")[0]})
		return
	} else {
		log.Logger.Infof("| %s | %s", trace, err)
		ctx.JSONP(http.StatusBadRequest, gin.H{"trace": trace, "msg": utils.M(ctx, "BadRequest"), "data": nil})
		return
	}
}

func GetUser(ctx *gin.Context) {
	trace := middleware.GetTraceID(ctx)
	user := middleware.GetUser(ctx)
	ctx.JSONP(http.StatusOK, gin.H{"trace": trace, "msg": utils.M(ctx, "Success"), "data": utils.TidyRetData(user, "password")[0]})
	return
}

func UpdateUser(ctx *gin.Context) {
	trace := middleware.GetTraceID(ctx)
	user := middleware.GetUser(ctx)
	var updateUserForm UpdateUserForm
	if err := ctx.ShouldBind(&updateUserForm); err == nil {
		data := utils.Form2Map(updateUserForm)
		if d, ok := data["password"]; ok {
			data["password"] = utils.HashPassword(d.(string))
		}
		if user.Type == "admin" {
			delete(data, "hidden")
			delete(data, "banned")
		}
		ok, msg := db.UpdateUser(user, data)
		if !ok {
			log.Logger.Infof("| %s | %s", trace, msg)
			ctx.JSONP(http.StatusInternalServerError, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
			return
		}
		// 避免回显数据时携带password hash
		delete(data, "password")
		ctx.JSONP(http.StatusOK, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": data})
		return
	} else {
		log.Logger.Infof("| %s | %s", trace, err)
		ctx.JSONP(http.StatusBadRequest, gin.H{"trace": trace, "msg": utils.M(ctx, "BadRequest"), "data": nil})
		return
	}
}

func DeleteUser(ctx *gin.Context) {
	trace := middleware.GetTraceID(ctx)
	user := middleware.GetUser(ctx)
	ok, msg := db.DeleteUser(user)
	if !ok {
		log.Logger.Infof("| %s | %s", trace, msg)
		ctx.JSONP(http.StatusInternalServerError, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
		return
	}
	ctx.JSONP(http.StatusOK, gin.H{"trace": trace, "msg": utils.M(ctx, "Success"), "data": nil})
	return
}

func CreateTeam(ctx *gin.Context) {
	var createTeamForm CreateTeamForm
	trace := middleware.GetTraceID(ctx)
	if err := ctx.ShouldBind(&createTeamForm); err == nil {
		teamName, contestID, userID := createTeamForm.Name, createTeamForm.ContestID, createTeamForm.UserID
		contest, ok, msg := db.GetContestByID(contestID)
		if !ok {
			log.Logger.Infof("| %s | %s", trace, msg)
			ctx.JSONP(http.StatusInternalServerError, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
			return
		}
		user, ok, msg := db.GetUserByID(userID)
		if !ok {
			log.Logger.Infof("| %s | %s", trace, msg)
			ctx.JSONP(http.StatusInternalServerError, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
			return
		}
		team, ok, msg := db.CreateTeam(contest, user, teamName)
		if !ok {
			log.Logger.Infof("| %s | %s", trace, msg)
			ctx.JSONP(http.StatusInternalServerError, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
			return
		}
		log.Logger.Infof("| %s | %s:%d create", trace, team.Name, team.ID)
		ctx.JSONP(http.StatusOK, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": utils.TidyRetData(team)[0]})
		return
	} else {
		log.Logger.Infof("| %s | %s", trace, err)
		ctx.JSONP(http.StatusBadRequest, gin.H{"trace": trace, "msg": utils.M(ctx, "BadRequest"), "data": nil})
		return
	}
}

func GetTeam(ctx *gin.Context) {
	trace := middleware.GetTraceID(ctx)
	team := middleware.GetTeam(ctx)
	ctx.JSONP(http.StatusOK, gin.H{"trace": trace, "msg": utils.M(ctx, "Success"), "data": utils.TidyRetData(team)[0]})
	return
}

func UpdateTeam(ctx *gin.Context) {
	trace := middleware.GetTraceID(ctx)
	team := middleware.GetTeam(ctx)
	var updateTeamForm UpdateTeamForm
	if err := ctx.ShouldBind(&updateTeamForm); err == nil {
		data := utils.Form2Map(updateTeamForm)
		if captcha, ok := data["captcha"]; ok && captcha == "refresh" {
			data["captcha"] = utils.RandomString()
		}
		ok, msg := db.UpdateTeam(team, data)
		if !ok {
			log.Logger.Infof("| %s | %s", trace, msg)
			ctx.JSONP(http.StatusInternalServerError, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
			return
		}
		ctx.JSONP(http.StatusOK, gin.H{"trace": trace, "msg": utils.M(ctx, "Success"), "data": data})
		return
	} else {
		log.Logger.Infof("| %s | %s", trace, err)
		ctx.JSONP(http.StatusBadRequest, gin.H{"trace": trace, "msg": utils.M(ctx, "BadRequest"), "data": nil})
		return
	}
}

func DeleteTeam(ctx *gin.Context) {
	trace := middleware.GetTraceID(ctx)
	team := middleware.GetTeam(ctx)
	ok, msg := db.DeleteTeam(team)
	if !ok {
		log.Logger.Infof("| %s | %s", trace, msg)
		ctx.JSONP(http.StatusInternalServerError, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
		return
	}
	ctx.JSONP(http.StatusOK, gin.H{"trace": trace, "msg": utils.M(ctx, "Success"), "data": nil})
	return
}

func CreateContest(ctx *gin.Context) {
	var createContestForm CreateContestForm
	trace := middleware.GetTraceID(ctx)
	if err := ctx.ShouldBind(&createContestForm); err == nil {
		contest, ok, msg := db.CreateContest(createContestForm.Name)
		if !ok {
			log.Logger.Infof("| %s | %s", trace, msg)
			ctx.JSONP(http.StatusInternalServerError, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
			return
		}
		log.Logger.Infof("| %s | %s:%d create", trace, contest.Name, contest.ID)
		ctx.JSONP(http.StatusOK, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": utils.TidyRetData(contest)[0]})
		return
	} else {
		log.Logger.Infof("| %s | %s", trace, err)
		ctx.JSONP(http.StatusBadRequest, gin.H{"trace": trace, "msg": utils.M(ctx, "BadRequest"), "data": nil})
		return
	}
}

func GetContest(ctx *gin.Context) {
	trace := middleware.GetTraceID(ctx)
	contest := middleware.GetContest(ctx)
	ctx.JSONP(http.StatusOK, gin.H{"trace": trace, "msg": utils.M(ctx, "Success"), "data": utils.TidyRetData(contest)[0]})
	return
}

func UpdateContest(ctx *gin.Context) {
	trace := middleware.GetTraceID(ctx)
	contest := middleware.GetContest(ctx)
	var updateContestForm UpdateContestForm
	if err := ctx.ShouldBind(&updateContestForm); err == nil {
		data := utils.Form2Map(updateContestForm)
		if duration, ok := data["duration"]; ok {
			delete(data, "duration")
			data["duration"] = time.Second * time.Duration(duration.(uint64))
		}
		if captcha, ok := data["captcha"]; ok && captcha == "refresh" {
			data["captcha"] = utils.RandomString()
		}
		if size, ok := data["size"]; ok && size.(uint) < contest.Size {
			msg := "TeamSizeIsLess"
			log.Logger.Infof("| %s | %s", trace, msg)
			ctx.JSONP(http.StatusBadRequest, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
			return
		}
		ok, msg := db.UpdateContest(contest, data)
		if !ok {
			log.Logger.Infof("| %s | %s", trace, msg)
			ctx.JSONP(http.StatusInternalServerError, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
			return
		}
		ctx.JSONP(http.StatusOK, gin.H{"trace": trace, "msg": utils.M(ctx, "Success"), "data": data})
		return
	} else {
		log.Logger.Infof("| %s | %s", trace, err)
		ctx.JSONP(http.StatusBadRequest, gin.H{"trace": trace, "msg": utils.M(ctx, "BadRequest"), "data": nil})
		return
	}
}

func DeleteContest(ctx *gin.Context) {
	trace := middleware.GetTraceID(ctx)
	contest := middleware.GetContest(ctx)
	ok, msg := db.DeleteContest(contest)
	if !ok {
		log.Logger.Infof("| %s | %s", trace, msg)
		ctx.JSONP(http.StatusInternalServerError, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
		return
	}
	ctx.JSONP(http.StatusOK, gin.H{"trace": trace, "msg": utils.M(ctx, "Success"), "data": nil})
	return
}

func GetFiles(ctx *gin.Context) {
	trace := middleware.GetTraceID(ctx)
	var getFilesForm GetFilesForm
	if err := ctx.ShouldBind(&getFilesForm); err == nil {
		files, total, ok, msg := db.GetFiles(getFilesForm.Limit, getFilesForm.Offset)
		if !ok {
			log.Logger.Infof("| %s | %s", trace, msg)
			ctx.JSONP(http.StatusInternalServerError, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
			return
		}
		ctx.JSONP(http.StatusOK, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": gin.H{
			"contests": utils.TidyRetData(files),
			"total":    total,
		}})
	} else {
		log.Logger.Infof("| %s | %s", trace, err)
		ctx.JSONP(http.StatusBadRequest, gin.H{"trace": trace, "msg": utils.M(ctx, "BadRequest"), "data": nil})
		return
	}
}

func DeleteFile(ctx *gin.Context) {
	trace := middleware.GetTraceID(ctx)
	type fileIDUri struct {
		FileID string `uri:"fileID" binding:"required"`
	}
	var fileID fileIDUri
	if err := ctx.ShouldBindUri(&fileID); err != nil {
		msg := "BadRequest"
		log.Logger.Infof("| %s | %s", trace, msg)
		ctx.JSONP(http.StatusBadRequest, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
		ctx.Abort()
		return
	}
	file, ok, msg := db.GetFile(fileID.FileID)
	if !ok {
		log.Logger.Infof("| %s | %s", trace, msg)
		ctx.JSONP(http.StatusNotFound, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": fileID.FileID})
		return
	}
	ok, msg = db.DeleteFile(fileID.FileID)
	if !ok {
		log.Logger.Infof("| %s | %s", trace, msg)
		ctx.JSONP(http.StatusInternalServerError, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": fileID.FileID})
		return
	}
	_ = os.Remove(file.Path)
	ctx.JSONP(http.StatusOK, gin.H{"trace": trace, "msg": utils.M(ctx, "Success"), "data": nil})
	return
}
