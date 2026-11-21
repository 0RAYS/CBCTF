package router

import (
	"sort"
	"time"

	"github.com/gin-gonic/gin"

	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
)

// formatNodeImages reports the cached detection result per node, limited to
// the target images. The node list stays live, so a node that joined after
// the detection is still selectable.
func formatNodeImages(nodes []string, detection *model.ImageDetection, targetImages []string) []gin.H {
	targetSet := make(map[string]struct{}, len(targetImages))
	for _, image := range targetImages {
		targetSet[image] = struct{}{}
	}
	results := make(map[string]model.NodeImageDetection)
	if detection != nil {
		for _, node := range detection.Nodes {
			results[node.Node] = node
		}
	}
	data := make([]gin.H, 0, len(nodes))
	for _, node := range nodes {
		current := make([]string, 0)
		missing := make([]string, 0)
		if result, ok := results[node]; ok {
			for _, image := range result.Present {
				if _, ok := targetSet[image]; ok {
					current = append(current, image)
				}
			}
			for _, image := range result.Absent {
				if _, ok := targetSet[image]; ok {
					missing = append(missing, image)
				}
			}
		}
		sort.Strings(current)
		sort.Strings(missing)
		data = append(data, gin.H{
			"node": node,
			// images holds the images the node holds, missing the ones it lacks.
			"images":  current,
			"missing": missing,
		})
	}
	sort.Slice(data, func(i, j int) bool { return data[i]["node"].(string) < data[j]["node"].(string) })
	return data
}

func detectedAt(detection *model.ImageDetection) string {
	if detection == nil || detection.DetectedAt.IsZero() {
		return ""
	}
	// The nanosecond part keeps two consecutive detections apart.
	return detection.DetectedAt.Format(time.RFC3339Nano)
}

func GetImages(ctx *gin.Context) {
	targetImages, ret := service.ListChallengeImages(db.DB.WithContext(ctx.Request.Context()))
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	nodes, ret := service.ListNodeNames(ctx.Request.Context())
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	detection := service.LoadImageDetection(ctx.Request.Context())
	resp.JSON(ctx, model.SuccessRetVal(gin.H{
		"nodes":         formatNodeImages(nodes, detection, targetImages),
		"target_images": targetImages,
		"detected_at":   detectedAt(detection),
	}))
}

func GetContestImage(ctx *gin.Context) {
	contest := middleware.GetContest(ctx)
	targetImages, ret := service.ListContestImages(db.DB.WithContext(ctx.Request.Context()), contest)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	nodes, ret := service.ListNodeNames(ctx.Request.Context())
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	detection := service.LoadImageDetection(ctx.Request.Context())
	resp.JSON(ctx, model.SuccessRetVal(gin.H{
		"nodes":         formatNodeImages(nodes, detection, targetImages),
		"target_images": targetImages,
		"detected_at":   detectedAt(detection),
	}))
}

// DetectImages starts a detection run for the whole challenge library.
func DetectImages(ctx *gin.Context) {
	targetImages, ret := service.ListChallengeImages(db.DB.WithContext(ctx.Request.Context()))
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ret = service.DetectImages(ctx.Request.Context(), targetImages)
	ctx.Set(middleware.CTXEventTypeKey, model.DetectImageEventType)
	ctx.Set(middleware.CTXEventSuccessKey, ret.OK)
	resp.JSON(ctx, ret)
}

// DetectContestImages starts a detection run for one contest.
func DetectContestImages(ctx *gin.Context) {
	contest := middleware.GetContest(ctx)
	targetImages, ret := service.ListContestImages(db.DB.WithContext(ctx.Request.Context()), contest)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ret = service.DetectImages(ctx.Request.Context(), targetImages)
	ctx.Set(middleware.CTXEventTypeKey, model.DetectImageEventType)
	ctx.Set(middleware.CTXEventSuccessKey, ret.OK)
	resp.JSON(ctx, ret)
}

func PullImages(ctx *gin.Context) {
	var form dto.PullImageForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.PullImageEventType)
	ret := service.PullContestImage(ctx.Request.Context(), form)
	ctx.Set(middleware.CTXEventSuccessKey, ret.OK)
	resp.JSON(ctx, ret)
}
