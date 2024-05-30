package localizer

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/text/language"
)

var bundle = i18n.NewBundle(language.English)
var localizer *i18n.Localizer = nil

func TranslateBundleSetup(path string, lang ...string) (*i18n.Bundle, error) {
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	_, err := bundle.LoadMessageFile(path)
	if err != nil {
		return nil, err
	}
	localizer = i18n.NewLocalizer(bundle, lang...)
	return bundle, nil
}

func GetLocalizer() *i18n.Localizer {
	if localizer == nil {
		localizer = i18n.NewLocalizer(bundle, "en-US", "en")
	}
	return localizer
}
