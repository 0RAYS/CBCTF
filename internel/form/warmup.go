package form

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/utils"
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"slices"
)

var allowedPullPolicy = []corev1.PullPolicy{corev1.PullAlways, corev1.PullNever, corev1.PullIfNotPresent}

type WarmUpImageForm struct {
	Images     []string `form:"images" json:"images" binding:"required"`
	PullPolicy string   `form:"pull_policy" json:"pull_policy" binding:"required"`
}

func (f *WarmUpImageForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		return false, i18n.BadRequest
	}
	f.PullPolicy = utils.ToTitle(f.PullPolicy)
	if !slices.Contains(allowedPullPolicy, corev1.PullPolicy(f.PullPolicy)) {
		f.PullPolicy = string(corev1.PullNever)
	}
	return true, i18n.Success
}
