package dto

import (
	"testing"

	"CBCTF/internal/i18n"
)

func TestUpdateTeamScoreFormRequiresAReasonAndANonZeroScore(t *testing.T) {
	form := UpdateTeamScoreForm{Score: 0, Reason: "  "}
	if ret := form.Validate(nil); ret.OK || ret.Msg != i18n.Model.Team.ScoreInvalid {
		t.Fatalf("zero score accepted: %+v", ret)
	}

	form = UpdateTeamScoreForm{Score: 100, Reason: "   "}
	if ret := form.Validate(nil); ret.OK || ret.Msg != i18n.Model.Team.ScoreReasonRequired {
		t.Fatalf("empty reason accepted: %+v", ret)
	}

	form = UpdateTeamScoreForm{Score: -50, Reason: " 题目分数配置错误 "}
	if ret := form.Validate(nil); !ret.OK {
		t.Fatalf("valid adjustment rejected: %+v", ret)
	}
	if form.Reason != "题目分数配置错误" {
		t.Fatalf("reason = %q, want the trimmed reason", form.Reason)
	}
}
