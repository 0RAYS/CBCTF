package xlsx

import (
	"bytes"
	"testing"
	"time"

	"github.com/xuri/excelize/v2"

	"CBCTF/internal/view"
)

func leaderboardFixture() view.LeaderboardExportView {
	solvedAt := time.Date(2026, 1, 2, 15, 4, 5, 0, time.UTC)
	return view.LeaderboardExportView{
		ContestName: "Demo",
		Challenges: []view.LeaderboardChallengeView{
			{ID: 10, Name: "Login", Category: "Web"},
			{ID: 11, Name: "Buffer", Category: "Pwn"},
		},
		Teams: []view.LeaderboardTeamView{
			{
				Rank: 1, Name: "Team A", Score: 1234.5, MemberCount: 3, SolvedCount: 1,
				LastSolve: solvedAt,
				Scores:    map[uint]float64{10: 1000},
				SolvedAt:  map[uint]time.Time{10: solvedAt},
			},
			{Rank: 2, Name: "Team B", Score: 0, MemberCount: 1, SolvedCount: 0},
		},
	}
}

func TestLeaderboardWorkbookHeaders(t *testing.T) {
	file, err := LeaderboardWorkbook(leaderboardFixture())
	if err != nil {
		t.Fatalf("build workbook: %v", err)
	}
	headers := []string{
		"排名 Rank", "队伍 Team", "总分 Score", "解题数 Solved", "成员数 Members", "最后解题 Last Solve",
		"Web-Login 得分 Score", "Web-Login 解题时间 Solved At",
		"Pwn-Buffer 得分 Score", "Pwn-Buffer 解题时间 Solved At",
	}
	for index, want := range headers {
		cell, err := excelize.CoordinatesToCellName(index+1, 1)
		if err != nil {
			t.Fatalf("cell name of column %d: %v", index+1, err)
		}
		got, err := file.GetCellValue(leaderboardSheetName, cell)
		if err != nil {
			t.Fatalf("read %s: %v", cell, err)
		}
		if got != want {
			t.Fatalf("header %s = %q, want %q", cell, got, want)
		}
	}
}

func TestLeaderboardWorkbookRows(t *testing.T) {
	file, err := LeaderboardWorkbook(leaderboardFixture())
	if err != nil {
		t.Fatalf("build workbook: %v", err)
	}
	cases := []struct {
		cell string
		want string
	}{
		{"A2", "1"},
		{"B2", "Team A"},
		{"C2", "1234.50"},
		{"D2", "1"},
		{"E2", "3"},
		{"G2", "1000.00"},
		{"H2", "2026-01-02 15:04:05"},
		{"I2", ""},
		{"J2", ""},
		{"B3", "Team B"},
		{"C3", "0.00"},
		{"G3", ""},
		{"H3", ""},
	}
	for _, testCase := range cases {
		got, err := file.GetCellValue(leaderboardSheetName, testCase.cell)
		if err != nil {
			t.Fatalf("read %s: %v", testCase.cell, err)
		}
		if got != testCase.want {
			t.Fatalf("cell %s = %q, want %q", testCase.cell, got, testCase.want)
		}
	}
}

func TestLeaderboardWorkbookIsReadableArchive(t *testing.T) {
	file, err := LeaderboardWorkbook(leaderboardFixture())
	if err != nil {
		t.Fatalf("build workbook: %v", err)
	}
	buffer, err := file.WriteToBuffer()
	if err != nil {
		t.Fatalf("encode workbook: %v", err)
	}
	if buffer.Len() == 0 {
		t.Fatal("encoded workbook is empty")
	}
	reopened, err := excelize.OpenReader(bytes.NewReader(buffer.Bytes()))
	if err != nil {
		t.Fatalf("open workbook: %v", err)
	}
	rows, err := reopened.GetRows(leaderboardSheetName)
	if err != nil {
		t.Fatalf("read rows: %v", err)
	}
	if len(rows) != 3 {
		t.Fatalf("row count = %d, want 3", len(rows))
	}
}
