package model

import "time"

// 证据事件的来源类型。
const (
	CheatSourceWebRequest     = "web_request"
	CheatSourceFlagSubmission = "flag_submission"
	CheatSourceFlagSolve      = "flag_solve"
	CheatSourceVictimAccess   = "victim_access"
)

// 证据涉及的角色类型。
const (
	CheatActorKindUser = "User"
	CheatActorKindTeam = "Team"
)

// CheatEvidence 是随线索一起保存的类型相关证据链。
// ReasonType 与线索的判定类型一致，Subject 是判定指向的具体对象：共享的 IP 或被提交的 flag。
// Actors 是判定涉及的用户或队伍，Timeline 是判定背后的观测事件，Metrics 概括观测窗口。
// 只有与判定类型对应的部分会被填充，其余部分为空。
type CheatEvidence struct {
	ReasonType CheatReasonType      `json:"reason_type"`
	Subject    string               `json:"subject"`
	Actors     []CheatActor         `json:"actors"`
	Timeline   []CheatTimelineItem  `json:"timeline"`
	Metrics    CheatMetrics         `json:"metrics"`
	Web        *CheatWebEvidence    `json:"web,omitempty"`
	Victim     *CheatVictimEvidence `json:"victim,omitempty"`
	Flag       *CheatFlagEvidence   `json:"flag,omitempty"`
}

// CheatActor 是判定涉及的一个用户或队伍。
type CheatActor struct {
	Kind      string    `json:"kind"`
	ID        uint      `json:"id"`
	Label     string    `json:"label,omitempty"`
	TeamIDs   []uint    `json:"team_ids,omitempty"`
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
	Count     int64     `json:"count"`
}

// CheatTimelineItem 是判定背后的一次观测事件。
type CheatTimelineItem struct {
	At        time.Time `json:"at"`
	Source    string    `json:"source"`
	ActorKind string    `json:"actor_kind"`
	ActorID   uint      `json:"actor_id"`
	Method    string    `json:"method,omitempty"`
	Path      string    `json:"path,omitempty"`
	Target    string    `json:"target,omitempty"`
	TargetID  uint      `json:"target_id,omitempty"`
	Status    int       `json:"status,omitempty"`
	LatencyMS int64     `json:"latency_ms,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
	Value     string    `json:"value,omitempty"`
	Detail    string    `json:"detail,omitempty"`
}

// CheatMetrics 概括一次判定的观测窗口。
type CheatMetrics struct {
	Events        int64     `json:"events"`
	ActorCount    int       `json:"actor_count"`
	FirstSeen     time.Time `json:"first_seen"`
	LastSeen      time.Time `json:"last_seen"`
	SpanSeconds   float64   `json:"span_seconds"`
	MinGapSeconds float64   `json:"min_gap_seconds"`
	MaxGapSeconds float64   `json:"max_gap_seconds"`
	Truncated     bool      `json:"truncated"`
}

// CheatWebEvidence 保存共享 IP 背后的网页请求、提交记录和涉及的题目。
type CheatWebEvidence struct {
	Requests    []CheatRequest        `json:"requests"`
	Submissions []CheatFlagSubmission `json:"submissions"`
	Challenges  []CheatChallengeRef   `json:"challenges"`
}

// CheatRequest 是共享 IP 上的一次网页请求。
type CheatRequest struct {
	At        time.Time `json:"at"`
	UserID    uint      `json:"user_id"`
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	Status    int       `json:"status"`
	LatencyMS int64     `json:"latency_ms"`
	UserAgent string    `json:"user_agent"`
}

// CheatChallengeRef 是判定涉及的题目副本及其使用者。
type CheatChallengeRef struct {
	ContestChallengeID uint   `json:"contest_challenge_id"`
	Name               string `json:"name"`
	ActorIDs           []uint `json:"actor_ids"`
}

// CheatFlagSubmission 是一次 flag 提交记录。
type CheatFlagSubmission struct {
	SubmissionID       uint      `json:"submission_id"`
	TeamID             uint      `json:"team_id"`
	UserID             uint      `json:"user_id"`
	ContestChallengeID uint      `json:"contest_challenge_id"`
	ChallengeName      string    `json:"challenge_name"`
	IP                 string    `json:"ip"`
	At                 time.Time `json:"at"`
	Value              string    `json:"value"`
	Solved             bool      `json:"solved"`
}

// CheatVictimEvidence 保存共享 IP 背后的靶机访问记录。
type CheatVictimEvidence struct {
	Victims  []CheatVictimRef    `json:"victims"`
	Accesses []CheatVictimAccess `json:"accesses"`
}

// CheatVictimRef 是被访问的靶机及其所属队伍与题目。
type CheatVictimRef struct {
	VictimID           uint      `json:"victim_id"`
	TeamID             uint      `json:"team_id"`
	ContestChallengeID uint      `json:"contest_challenge_id"`
	ChallengeName      string    `json:"challenge_name"`
	AccessCount        int64     `json:"access_count"`
	FirstSeen          time.Time `json:"first_seen"`
	LastSeen           time.Time `json:"last_seen"`
}

// CheatVictimAccess 是一次靶机访问观测。
type CheatVictimAccess struct {
	At       time.Time `json:"at"`
	TeamID   uint      `json:"team_id"`
	VictimID uint      `json:"victim_id"`
	Source   string    `json:"source"`
	Capture  string    `json:"capture"`
}

// CheatFlagEvidence 保存误用他人 flag 的提交证据。
type CheatFlagEvidence struct {
	ContestChallengeID uint                  `json:"contest_challenge_id"`
	ChallengeName      string                `json:"challenge_name"`
	SubmittingTeamID   uint                  `json:"submitting_team_id"`
	MatchedTeamIDs     []uint                `json:"matched_team_ids"`
	MatchedValues      []string              `json:"matched_values"`
	Owners             []CheatFlagOwner      `json:"owners"`
	Submissions        []CheatFlagSubmission `json:"submissions"`
}

// CheatFlagOwner 是被冒用的 flag 的合法持有者及其解题时间。
type CheatFlagOwner struct {
	TeamID     uint      `json:"team_id"`
	UserID     uint      `json:"user_id"`
	Value      string    `json:"value"`
	ObtainedAt time.Time `json:"obtained_at"`
}

// CheatEvidenceSummary 是列表视图使用的紧凑证据摘要。
type CheatEvidenceSummary struct {
	Subject       string          `json:"subject"`
	ReasonType    CheatReasonType `json:"reason_type"`
	Actors        int             `json:"actors"`
	Events        int64           `json:"events"`
	FirstSeen     time.Time       `json:"first_seen"`
	LastSeen      time.Time       `json:"last_seen"`
	SpanSeconds   float64         `json:"span_seconds"`
	MinGapSeconds float64         `json:"min_gap_seconds"`
	Truncated     bool            `json:"truncated"`
}

// Summary 返回证据链的紧凑摘要。
func (e CheatEvidence) Summary() CheatEvidenceSummary {
	return CheatEvidenceSummary{
		Subject:       e.Subject,
		ReasonType:    e.ReasonType,
		Actors:        len(e.Actors),
		Events:        e.Metrics.Events,
		FirstSeen:     e.Metrics.FirstSeen,
		LastSeen:      e.Metrics.LastSeen,
		SpanSeconds:   e.Metrics.SpanSeconds,
		MinGapSeconds: e.Metrics.MinGapSeconds,
		Truncated:     e.Metrics.Truncated,
	}
}
