package view

import "CBCTF/internal/model"

type ChallengeFlagView struct {
	ID    uint
	Value string
}

type ChallengeView struct {
	Challenge     model.Challenge
	Flags         []ChallengeFlagView
	DockerCompose string
	FileName      string
}

type SimpleChallengeView struct {
	Challenge model.Challenge
}
