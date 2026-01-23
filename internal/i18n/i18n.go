package i18n

import (
	_ "embed"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.yaml.in/yaml/v3"
	"golang.org/x/text/language"
)

var (
	Bundle *i18n.Bundle

	//go:embed locales/zh-CN.yaml
	cn []byte
	//go:embed locales/en.yaml
	en []byte
	//go:embed locales/code.yaml
	code []byte
)

func Init() {
	Bundle = i18n.NewBundle(language.English)
	Bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	Bundle.MustParseMessageFileBytes(cn, "zh-CN.yaml")
	Bundle.MustParseMessageFileBytes(en, "en.yaml")
	Bundle.MustParseMessageFileBytes(code, "und.yaml")
}

func GetLocalizer(lang string) *i18n.Localizer {
	return i18n.NewLocalizer(Bundle, lang)
}

func Translate(lang, key string, args ...map[string]any) string {
	config := i18n.LocalizeConfig{MessageID: key}
	if len(args) > 0 {
		config.TemplateData = args[0]
	}
	msg, err := GetLocalizer(lang).Localize(&config)
	if err != nil {
		return key
	}
	return msg
}
