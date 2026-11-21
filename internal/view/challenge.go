package view

import "CBCTF/internal/model"

type ChallengeFlagView struct {
	ID    uint
	Value string
}

type ChallengeView struct {
	Unavailable   []string
	Challenge     model.Challenge
	ContestName   string
	Flags         []ChallengeFlagView
	DockerCompose string
	FileName      string
}
