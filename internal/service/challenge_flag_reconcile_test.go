package service

import (
	"testing"

	"CBCTF/internal/model"
)

func TestMatchChallengeFlagsKeepsFlagIdentityWhenValueChanges(t *testing.T) {
	existing := []model.ChallengeFlag{
		{ID: 1, Value: "flag{static:aaa}", Binding: model.FlagBinding{PodKey: "web", ContainerKey: "web", Type: model.EnvFlagBindingType, Target: "FLAG"}},
		{ID: 2, Value: "flag{static:bbb}", Binding: model.FlagBinding{PodKey: "web", ContainerKey: "web", Type: model.EnvFlagBindingType, Target: "SECRET"}},
	}
	derived := []model.ChallengeFlag{
		{Value: "flag{static:ccc}", Binding: model.FlagBinding{PodKey: "web", ContainerKey: "web", Type: model.EnvFlagBindingType, Target: "FLAG"}},
		{Value: "flag{static:bbb}", Binding: model.FlagBinding{PodKey: "web", ContainerKey: "web", Type: model.EnvFlagBindingType, Target: "SECRET"}},
	}
	matches, removed := matchChallengeFlags(existing, derived)
	if len(matches) != 2 {
		t.Fatalf("matches=%d", len(matches))
	}
	if matches[0].ExistingID != 1 || matches[0].Value != "flag{static:ccc}" || !matches[0].Changed {
		t.Fatalf("value change must keep the flag id: %+v", matches[0])
	}
	if matches[1].ExistingID != 2 || matches[1].Changed {
		t.Fatalf("unchanged flag must be reused: %+v", matches[1])
	}
	if len(removed) != 0 {
		t.Fatalf("removed=%v", removed)
	}
}

func TestMatchChallengeFlagsCreatesAndRemovesFlags(t *testing.T) {
	existing := []model.ChallengeFlag{
		{ID: 1, Value: "flag{static:aaa}", Binding: model.FlagBinding{Target: "FLAG"}},
		{ID: 2, Value: "flag{static:bbb}", Binding: model.FlagBinding{Target: "OLD"}},
	}
	derived := []model.ChallengeFlag{
		{Value: "flag{static:aaa}", Binding: model.FlagBinding{Target: "FLAG"}},
		{Value: "flag{static:ccc}", Binding: model.FlagBinding{Target: "NEW"}},
	}
	matches, removed := matchChallengeFlags(existing, derived)
	if len(matches) != 2 {
		t.Fatalf("matches=%d", len(matches))
	}
	if matches[0].ExistingID != 1 || matches[0].Changed {
		t.Fatalf("existing flag must be reused: %+v", matches[0])
	}
	if matches[1].ExistingID != 0 || matches[1].Value != "flag{static:ccc}" {
		t.Fatalf("new flag must be created: %+v", matches[1])
	}
	if len(removed) != 1 || removed[0] != 2 {
		t.Fatalf("removed=%v", removed)
	}
}

func TestMatchChallengeFlagsKeepsMultipleFlagsOnTheSameBinding(t *testing.T) {
	existing := []model.ChallengeFlag{
		{ID: 1, Value: "flag{static:aaa}", Binding: model.FlagBinding{Target: "FILE"}},
		{ID: 2, Value: "flag{static:bbb}", Binding: model.FlagBinding{Target: "FILE"}},
	}
	derived := []model.ChallengeFlag{
		{Value: "flag{static:aaa}", Binding: model.FlagBinding{Target: "FILE"}},
		{Value: "flag{static:bbb}", Binding: model.FlagBinding{Target: "FILE"}},
	}
	matches, removed := matchChallengeFlags(existing, derived)
	if len(matches) != 2 || len(removed) != 0 {
		t.Fatalf("matches=%+v removed=%v", matches, removed)
	}
	if matches[0].ExistingID != 1 || matches[1].ExistingID != 2 {
		t.Fatalf("flags on the same binding must keep their identity: %+v", matches)
	}
}
