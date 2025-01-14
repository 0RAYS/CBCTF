package i18n

import (
	"embed"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
	"io/fs"
)

//go:embed locales/*
var localesFS embed.FS

var Bundle *i18n.Bundle

func Init() {
	Bundle = i18n.NewBundle(language.English)
	Bundle.RegisterUnmarshalFunc("yaml", yamlUnmarshalFunc)
	loadEmbeddedTranslations()
	//acceptLang := c.GetHeader("Accept-Language")
	//localizedMessage := i18n.NewLocalizer(bundle, acceptLang).MustLocalize(&i18n.LocalizeConfig{
	//	MessageID: "hello_world",
	//})
}

func loadEmbeddedTranslations() {
	files, err := fs.ReadDir(localesFS, "locales")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		content, err := localesFS.ReadFile("locales/" + file.Name())
		if err != nil {
			panic(err)
		}
		Bundle.MustParseMessageFileBytes(content, "locales/"+file.Name())
	}
}

func yamlUnmarshalFunc(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}
