package middleware

import (
	"CBCTF/internal/log"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strings"
)

var resp = map[string]map[string]interface{}{
	"Success":               {"zh-CN": "操作成功", "en-US": "Success", "code": 200},
	"BadRequest":            {"zh-CN": "请求错误", "en-US": "Bad request", "code": 400},
	"CreateUserError":       {"zh-CN": "创建用户失败", "en-US": "Failed to create user", "code": 500},
	"CreateTeamError":       {"zh-CN": "创建队伍失败", "en-US": "Failed to create team", "code": 500},
	"CreateContestError":    {"zh-CN": "创建赛事失败", "en-US": "Failed to create contest", "code": 500},
	"CreateFileRecordError": {"zh-CN": "保存文件失败", "en-US": "Failed to save file", "code": 500},
	"DeleteContestError":    {"zh-CN": "删除赛事失败", "en-US": "Failed to delete contest", "code": 500},
	"DeleteFileError":       {"zh-CN": "删除文件失败", "en-US": "Failed to delete file", "code": 500},
	"DeleteTeamError":       {"zh-CN": "删除队伍失败", "en-US": "Failed to delete team", "code": 500},
	"DeleteUserError":       {"zh-CN": "删除用户失败", "en-US": "Failed to delete user", "code": 500},
	"EmailExists":           {"zh-CN": "该邮箱已注册", "en-US": "Email already registered", "code": 400},
	"FileNotAllowed":        {"zh-CN": "不支持的文件类型", "en-US": "Unsupported file type", "code": 400},
	"FileNotFound":          {"zh-CN": "文件不存在", "en-US": "File not found", "code": 404},
	"Forbidden":             {"zh-CN": "禁止访问", "en-US": "Forbidden", "code": 403},
	"IsNotPlayer":           {"zh-CN": "该用户未加入任何队伍", "en-US": "User has not joined any team", "code": 400},
	"InvalidEmail":          {"zh-CN": "邮箱地址无效", "en-US": "Invalid Email address", "code": 400},
	"JoinTeamError":         {"zh-CN": "加入队伍失败", "en-US": "Failed to join team", "code": 500},
	"LeaveTeamError":        {"zh-CN": "退出队伍失败", "en-US": "Failed to leave team", "code": 500},
	"NameOrPasswordError":   {"zh-CN": "用户名或密码错误，请重试", "en-US": "Username or password error, please try again", "code": 401},
	"PasswordError":         {"zh-CN": "密码错误", "en-US": "Incorrect password", "code": 401},
	"PasswordSame":          {"zh-CN": "新密码与旧密码相同", "en-US": "New password is the same as the old one", "code": 400},
	"RepeatPlayer":          {"zh-CN": "该用户已在其他队伍中", "en-US": "User is already in another team", "code": 400},
	"TeamIsFull":            {"zh-CN": "队伍人数已满", "en-US": "Team is full", "code": 400},
	"TeamNameExists":        {"zh-CN": "队伍名已被占用", "en-US": "Team name already taken", "code": 400},
	"TeamNotFound":          {"zh-CN": "队伍不存在", "en-US": "Team not found", "code": 404},
	"TeamFull":              {"zh-CN": "队伍已满", "en-US": "Team is full", "code": 400},
	"Unauthorized":          {"zh-CN": "未登录", "en-US": "Unauthorized", "code": 401},
	"UnknownError":          {"zh-CN": "未知错误，请联系管理员", "en-US": "UnknownError, please contact the administrator", "code": 500},
	"UnverifiedEmail":       {"zh-CN": "邮箱未验证", "en-US": "Email not verified", "code": 403},
	"UploadFileError":       {"zh-CN": "文件上传失败", "en-US": "File upload failed", "code": 500},
	"UserNameExists":        {"zh-CN": "用户名已被占用", "en-US": "Username already taken", "code": 400},
	"UserNotFound":          {"zh-CN": "用户不存在", "en-US": "User not found", "code": 404},
	"UserNotInTeam":         {"zh-CN": "用户不在队伍中", "en-US": "User is not in the team", "code": 400},
	"ContestNameExists":     {"zh-CN": "赛事名已被占用", "en-US": "Contest name already taken", "code": 400},
	"ContestNotFound":       {"zh-CN": "赛事不存在", "en-US": "Contest not found", "code": 404},
}

type i18nResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *i18nResponseWriter) Write(p []byte) (n int, err error) {
	return w.body.Write(p)
}

type Data struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Data  any    `json:"data"`
	Trace string `json:"trace"`
}

func I18n() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		w := &i18nResponseWriter{
			ResponseWriter: ctx.Writer,
			body:           bytes.NewBufferString(""),
		}
		ctx.Writer = w

		ctx.Next()

		var result Data
		old := w.body.String()

		err := json.Unmarshal([]byte(old), &result)
		if err != nil {
			_, _ = w.ResponseWriter.Write([]byte(old))
			return
		}
		result.Code = resp[result.Msg]["code"].(int)
		language := ctx.GetHeader("Accept-Language")
		if strings.HasPrefix(language, "zh-CN") {
			language = "zh-CN"
		} else {
			language = "en-US"
		}
		result.Msg = resp[result.Msg][language].(string)
		result.Trace = GetTraceID(ctx)
		ret, err := json.Marshal(result)
		if err != nil {
			log.Logger.Errorf("Rewrite response error: %v", err)
			return
		}
		defer w.body.Reset()
		_, _ = w.ResponseWriter.Write(ret)
	}
}
