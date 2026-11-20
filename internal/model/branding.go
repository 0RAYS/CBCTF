package model

import (
	"fmt"

	"CBCTF/internal/config"
)

const DefaultBrandingCode = "default"

type LocalizedText struct {
	ZhCN string `json:"zh_cn"`
	En   string `json:"en"`
}

type BrandingHeroContent struct {
	TitlePrefix     LocalizedText `json:"title_prefix"`
	TitleHighlight  LocalizedText `json:"title_highlight"`
	TitleSuffix     LocalizedText `json:"title_suffix"`
	Subtitle        LocalizedText `json:"subtitle"`
	PrimaryAction   LocalizedText `json:"primary_action"`
	SecondaryAction LocalizedText `json:"secondary_action"`
}

type BrandingSectionContent struct {
	TitlePrefix    LocalizedText `json:"title_prefix"`
	TitleHighlight LocalizedText `json:"title_highlight"`
	Subtitle       LocalizedText `json:"subtitle"`
}

type BrandingActionSectionContent struct {
	TitlePrefix    LocalizedText `json:"title_prefix"`
	TitleHighlight LocalizedText `json:"title_highlight"`
	Subtitle       LocalizedText `json:"subtitle"`
	Action         LocalizedText `json:"action"`
}

type BrandingHomeContent struct {
	Hero           BrandingHeroContent          `json:"hero"`
	ChallengeTypes BrandingSectionContent       `json:"challenge_types"`
	Upcoming       BrandingActionSectionContent `json:"upcoming"`
	Leaderboard    BrandingActionSectionContent `json:"leaderboard"`
}

type Branding struct {
	Code               string              `gorm:"type:varchar(64);uniqueIndex:idx_brandings_code_active,where:deleted_at IS NULL;not null" json:"code"`
	SiteName           LocalizedText       `gorm:"serializer:json;type:jsonb;not null;default:'{}'" json:"site_name"`
	AdminName          LocalizedText       `gorm:"serializer:json;type:jsonb;not null;default:'{}'" json:"admin_name"`
	BrowserTitle       LocalizedText       `gorm:"serializer:json;type:jsonb;not null;default:'{}'" json:"browser_title"`
	BrowserDescription LocalizedText       `gorm:"serializer:json;type:jsonb;not null;default:'{}'" json:"browser_description"`
	FooterCopyright    LocalizedText       `gorm:"serializer:json;type:jsonb;not null;default:'{}'" json:"footer_copyright"`
	FooterICPNumber    string              `gorm:"type:varchar(128)" json:"footer_icp_number"`
	FooterICPLink      string              `gorm:"type:varchar(255)" json:"footer_icp_link"`
	FooterContactEmail string              `gorm:"type:varchar(255)" json:"footer_contact_email"`
	FooterGithubURL    string              `gorm:"type:varchar(255)" json:"footer_github_url"`
	HomeLogo           FileURL             `json:"home_logo"`
	HomeLogoAlt        LocalizedText       `gorm:"serializer:json;type:jsonb;not null;default:'{}'" json:"home_logo_alt"`
	Home               BrandingHomeContent `gorm:"serializer:json;type:jsonb;not null;default:'{}'" json:"home"`
	BaseModel
}

func DefaultBranding() Branding {
	return Branding{
		Code: DefaultBrandingCode,
		SiteName: LocalizedText{
			ZhCN: "CBCTF",
			En:   "CBCTF",
		},
		AdminName: LocalizedText{
			ZhCN: "CBCTF",
			En:   "CBCTF",
		},
		BrowserTitle: LocalizedText{
			ZhCN: "CBCTF",
			En:   "CBCTF",
		},
		BrowserDescription: LocalizedText{
			ZhCN: "CBCTF competition platform",
			En:   "CBCTF competition platform",
		},
		FooterCopyright: LocalizedText{
			ZhCN: "© 2026 CBCTF",
			En:   "© 2026 CBCTF",
		},
		FooterICPNumber:    "浙ICP备12028388号-2",
		FooterICPLink:      "https://beian.miit.gov.cn/",
		FooterContactEmail: "support@0rays.club",
		FooterGithubURL:    "https://github.com/0RAYS/CBCTF",
		HomeLogo:           FileURL(fmt.Sprintf("%s/platform/logo.png", config.Env.Host)),
		HomeLogoAlt: LocalizedText{
			ZhCN: "CBCTF Logo",
			En:   "CBCTF logo",
		},
		Home: BrandingHomeContent{
			Hero: BrandingHeroContent{
				TitlePrefix: LocalizedText{
					ZhCN: "Dive Deep into the",
					En:   "Dive Deep into the",
				},
				TitleHighlight: LocalizedText{
					ZhCN: "Cyber Security",
					En:   "Cyber Security",
				},
				TitleSuffix: LocalizedText{
					ZhCN: "Challenge",
					En:   "Challenge",
				},
				Subtitle: LocalizedText{
					ZhCN: "Join the elite community of hackers, compete in real-world challenges, and master the art of cybersecurity through hands-on experience.",
					En:   "Join the elite community of hackers, compete in real-world challenges, and master the art of cybersecurity through hands-on experience.",
				},
				PrimaryAction: LocalizedText{
					ZhCN: "START HACKING",
					En:   "START HACKING",
				},
				SecondaryAction: LocalizedText{
					ZhCN: "LEARN MORE",
					En:   "LEARN MORE",
				},
			},
			ChallengeTypes: BrandingSectionContent{
				TitlePrefix: LocalizedText{
					ZhCN: "Master All Aspects of",
					En:   "Master All Aspects of",
				},
				TitleHighlight: LocalizedText{
					ZhCN: "Cyber Security",
					En:   "Cyber Security",
				},
				Subtitle: LocalizedText{
					ZhCN: "From web exploitation to binary analysis, our challenges cover every major domain of cybersecurity",
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
					En:   "VIEW LEADERBOARD",
				},
			},
		},
	}
}
