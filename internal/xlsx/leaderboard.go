package xlsx

import (
	"fmt"

	"github.com/xuri/excelize/v2"

	"CBCTF/internal/view"
)

const (
	leaderboardSheetName = "Leaderboard"
	scoreFormat          = "0.00"
	timeFormat           = "yyyy-mm-dd hh:mm:ss"
)

// LeaderboardWorkbook 生成排行榜工作簿：一行一支队伍，固定列给出排名与总分，
// 每道比赛题目占两列，分别记录该队在这道题上的得分与首次解题时间。
func LeaderboardWorkbook(data view.LeaderboardExportView) (*excelize.File, error) {
	file := excelize.NewFile()
	sheet := leaderboardSheetName
	if err := file.SetSheetName(file.GetSheetName(0), sheet); err != nil {
		return nil, err
	}

	headers := []string{"排名 Rank", "队伍 Team", "总分 Score", "解题数 Solved", "成员数 Members", "最后解题 Last Solve"}
	labels := make([]string, 0, len(data.Challenges))
	occurrences := make(map[string]int, len(data.Challenges))
	for _, challenge := range data.Challenges {
		label := challenge.Name
		if challenge.Category != "" {
			label = challenge.Category + "-" + label
		}
		occurrences[label]++
		if occurrences[label] > 1 {
			label = fmt.Sprintf("%s #%d", label, occurrences[label])
		}
		labels = append(labels, label)
		headers = append(headers, label+" 得分 Score", label+" 解题时间 Solved At")
	}
	for index, header := range headers {
		cell, err := excelize.CoordinatesToCellName(index+1, 1)
		if err != nil {
			return nil, err
		}
		if err = file.SetCellStr(sheet, cell, header); err != nil {
			return nil, err
		}
	}

	for index, team := range data.Teams {
		row := []any{team.Rank, team.Name, team.Score, team.SolvedCount, team.MemberCount}
		if team.LastSolve.IsZero() {
			row = append(row, nil)
		} else {
			row = append(row, team.LastSolve)
		}
		for _, challenge := range data.Challenges {
			if score, ok := team.Scores[challenge.ID]; ok {
				row = append(row, score)
			} else {
				row = append(row, nil)
			}
			if solvedAt, ok := team.SolvedAt[challenge.ID]; ok {
				row = append(row, solvedAt)
			} else {
				row = append(row, nil)
			}
		}
		cell, err := excelize.CoordinatesToCellName(1, index+2)
		if err != nil {
			return nil, err
		}
		if err = file.SetSheetRow(sheet, cell, &row); err != nil {
			return nil, err
		}
	}

	if err := styleLeaderboard(file, sheet, headers, labels, len(data.Teams)); err != nil {
		return nil, err
	}
	return file, nil
}

func styleLeaderboard(file *excelize.File, sheet string, headers, labels []string, teamCount int) error {
	headerStyle, err := file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#334155"}},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
	})
	if err != nil {
		return err
	}
	scoreNumFmt := scoreFormat
	timeNumFmt := timeFormat
	scoreStyle, err := file.NewStyle(&excelize.Style{CustomNumFmt: &scoreNumFmt})
	if err != nil {
		return err
	}
	timeStyle, err := file.NewStyle(&excelize.Style{CustomNumFmt: &timeNumFmt})
	if err != nil {
		return err
	}

	lastColumn, err := excelize.ColumnNumberToName(len(headers))
	if err != nil {
		return err
	}
	lastRow := teamCount + 1
	if err = file.SetCellStyle(sheet, "A1", lastColumn+"1", headerStyle); err != nil {
		return err
	}
	if err = file.SetRowHeight(sheet, 1, 32); err != nil {
		return err
	}

	widths := map[string]float64{"A": 8, "B": 24, "C": 12, "D": 10, "E": 10, "F": 20}
	for column, width := range widths {
		if err = file.SetColWidth(sheet, column, column, width); err != nil {
			return err
		}
	}
	for index := range labels {
		scoreColumn, err := excelize.ColumnNumberToName(7 + index*2)
		if err != nil {
			return err
		}
		timeColumn, err := excelize.ColumnNumberToName(8 + index*2)
		if err != nil {
			return err
		}
		if err = file.SetColWidth(sheet, scoreColumn, scoreColumn, 14); err != nil {
			return err
		}
		if err = file.SetColWidth(sheet, timeColumn, timeColumn, 20); err != nil {
			return err
		}
		if lastRow > 1 {
			if err = file.SetCellStyle(sheet, scoreColumn+"2", scoreColumn+fmt.Sprint(lastRow), scoreStyle); err != nil {
				return err
			}
			if err = file.SetCellStyle(sheet, timeColumn+"2", timeColumn+fmt.Sprint(lastRow), timeStyle); err != nil {
				return err
			}
		}
	}
	if lastRow > 1 {
		if err = file.SetCellStyle(sheet, "C2", "C"+fmt.Sprint(lastRow), scoreStyle); err != nil {
			return err
		}
		if err = file.SetCellStyle(sheet, "F2", "F"+fmt.Sprint(lastRow), timeStyle); err != nil {
			return err
		}
	}

	if err = file.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		XSplit:      2,
		YSplit:      1,
		TopLeftCell: "C2",
		ActivePane:  "bottomRight",
	}); err != nil {
		return err
	}
	return file.AutoFilter(sheet, "A1:"+lastColumn+fmt.Sprint(lastRow), nil)
}
