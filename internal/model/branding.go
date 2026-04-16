package model

import (
	"CBCTF/internal/config"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

const DefaultBrandingCode = "default"

type LocalizedText struct {
	ZhCN string `json:"zh_cn"`
	En   string `json:"en"`
}

func (l LocalizedText) Value() (driver.Value, error) {
	return json.Marshal(l)
}

func (l *LocalizedText) Scan(value any) error {
	if err := scanJSON(value, l); err != nil {
		return fmt.Errorf("failed to scan LocalizedText value")
	}
	return nil
}

type BrandingHeroContent struct {
	TitlePrefix     LocalizedText `json:"title_prefix"`
	TitleHighlight  LocalizedText `json:"title_highlight"`
	TitleSuffix     LocalizedText `json:"title_suffix"`
	Subtitle        LocalizedText `json:"subtitle"`
	PrimaryAction   LocalizedText `json:"primary_action"`
	SecondaryAction LocalizedText `json:"secondary_action"`
}

func (b BrandingHeroContent) Value() (driver.Value, error) {
	return json.Marshal(b)
}

func (b *BrandingHeroContent) Scan(value any) error {
	if err := scanJSON(value, b); err != nil {
		return fmt.Errorf("failed to scan BrandingHeroContent value")
	}
	return nil
}

type BrandingSectionContent struct {
	TitlePrefix    LocalizedText `json:"title_prefix"`
	TitleHighlight LocalizedText `json:"title_highlight"`
	Subtitle       LocalizedText `json:"subtitle"`
}

func (b BrandingSectionContent) Value() (driver.Value, error) {
	return json.Marshal(b)
}

func (b *BrandingSectionContent) Scan(value any) error {
	if err := scanJSON(value, b); err != nil {
		return fmt.Errorf("failed to scan BrandingSectionContent value")
	}
	return nil
}

type BrandingActionSectionContent struct {
	TitlePrefix    LocalizedText `json:"title_prefix"`
	TitleHighlight LocalizedText `json:"title_highlight"`
	Subtitle       LocalizedText `json:"subtitle"`
	Action         LocalizedText `json:"action"`
}

func (b BrandingActionSectionContent) Value() (driver.Value, error) {
	return json.Marshal(b)
}

func (b *BrandingActionSectionContent) Scan(value any) error {
	if err := scanJSON(value, b); err != nil {
		return fmt.Errorf("failed to scan BrandingActionSectionContent value")
	}
	return nil
}

type BrandingHomeContent struct {
	Hero           BrandingHeroContent          `json:"hero"`
	ChallengeTypes BrandingSectionContent       `json:"challenge_types"`
	Upcoming       BrandingActionSectionContent `json:"upcoming"`
	Leaderboard    BrandingActionSectionContent `json:"leaderboard"`
}

func (b BrandingHomeContent) Value() (driver.Value, error) {
	return json.Marshal(b)
}

func (b *BrandingHomeContent) Scan(value any) error {
	if err := scanJSON(value, b); err != nil {
		return fmt.Errorf("failed to scan BrandingHomeContent value")
	}
	return nil
}

type Branding struct {
	Code               string              `gorm:"type:varchar(64);uniqueIndex:idx_brandings_code_active,where:deleted_at IS NULL;not null" json:"code"`
	SiteName           LocalizedText       `gorm:"type:jsonb;not null" json:"site_name"`
	AdminName          LocalizedText       `gorm:"type:jsonb;not null" json:"admin_name"`
	BrowserTitle       LocalizedText       `gorm:"type:jsonb;not null" json:"browser_title"`
	BrowserDescription LocalizedText       `gorm:"type:jsonb;not null" json:"browser_description"`
	FooterCopyright    LocalizedText       `gorm:"type:jsonb;not null" json:"footer_copyright"`
	HomeLogo           FileURL             `json:"home_logo"`
	HomeLogoAlt        LocalizedText       `gorm:"type:jsonb;not null" json:"home_logo_alt"`
	Home               BrandingHomeContent `gorm:"type:jsonb;not null" json:"home"`
	BaseModel
}

func DefaultBranding() Branding {
	return Branding{
		Code: DefaultBrandingCode,
		SiteName: LocalizedText{
			ZhCN: "深潜 CTF",
			En:   "DEEP DIVE CTF",
		},
		AdminName: LocalizedText{
			ZhCN: "深潜管理台",
			En:   "DEEP DIVE Admin",
		},
		BrowserTitle: LocalizedText{
			ZhCN: "深潜 CTF",
			En:   "DEEP DIVE CTF",
		},
		BrowserDescription: LocalizedText{
			ZhCN: "深潜 CTF 网络安全竞赛平台",
			En:   "DEEP DIVE CTF competition platform",
		},
		FooterCopyright: LocalizedText{
			ZhCN: "© 2025 深潜 CTF",
			En:   "© 2025 DEEP DIVE CTF",
		},
		HomeLogo: FileURL(fmt.Sprintf("%s/platform/logo.png", config.Env.Host)),
		HomeLogoAlt: LocalizedText{
			ZhCN: "深潜 CTF 首页 Logo",
			En:   "DEEP DIVE CTF home logo",
		},
		Home: BrandingHomeContent{
			Hero: BrandingHeroContent{
				TitlePrefix: LocalizedText{
					ZhCN: "深入探索",
					En:   "Dive Deep into the",
				},
				TitleHighlight: LocalizedText{
					ZhCN: "网络安全",
					En:   "Cyber Security",
				},
				TitleSuffix: LocalizedText{
					ZhCN: "挑战",
					En:   "Challenge",
				},
				Subtitle: LocalizedText{
					ZhCN: "加入深潜 CTF 社区, 在真实场景中对抗、练习与成长",
					En:   "Join the elite community of hackers, compete in real-world challenges, and master the art of cybersecurity through hands-on experience.",
				},
				PrimaryAction: LocalizedText{
					ZhCN: "立即参赛",
					En:   "START HACKING",
				},
				SecondaryAction: LocalizedText{
					ZhCN: "了解更多",
					En:   "LEARN MORE",
				},
			},
			ChallengeTypes: BrandingSectionContent{
				TitlePrefix: LocalizedText{
					ZhCN: "掌握",
					En:   "Master All Aspects of",
				},
				TitleHighlight: LocalizedText{
					ZhCN: "网络安全",
					En:   "Cyber Security",
				},
				Subtitle: LocalizedText{
					ZhCN: "从 Web 渗透到二进制分析, 覆盖网络安全核心方向",
					En:   "From web exploitation to binary analysis, our challenges cover every major domain of cybersecurity",
				},
			},
			Upcoming: BrandingActionSectionContent{
				TitlePrefix: LocalizedText{
					ZhCN: "近期",
					En:   "Upcoming",
				},
				TitleHighlight: LocalizedText{
					ZhCN: "赛事",
					En:   "Competitions",
				},
				Subtitle: LocalizedText{
					ZhCN: "报名即将开启的 CTF 赛事, 与高手同场竞技",
					En:   "Register now for our upcoming CTF events and compete with the best",
				},
				Action: LocalizedText{
					ZhCN: "查看全部赛事",
					En:   "VIEW ALL CONTESTS",
				},
			},
			Leaderboard: BrandingActionSectionContent{
				TitlePrefix: LocalizedText{
					ZhCN: "顶尖",
					En:   "Top",
				},
				TitleHighlight: LocalizedText{
					ZhCN: "战队",
					En:   "Performers",
				},
				Subtitle: LocalizedText{
					ZhCN: "查看当前积分榜上的领先队伍",
					En:   "Meet the elite teams leading our global leaderboard",
				},
				Action: LocalizedText{
					ZhCN: "查看排行榜",
					En:   "VIEW SCOREBOARD",
				},
			},
		},
	}
}
