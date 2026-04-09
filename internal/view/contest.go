package view

import "CBCTF/internal/model"

type ContestView struct {
	Contest     model.Contest
	TeamCount   int64
	UserCount   int64
	NoticeCount int64
	Highest     float64
	SolvedCount int64
	StatsReady  bool
}
