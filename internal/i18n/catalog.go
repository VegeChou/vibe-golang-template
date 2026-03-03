package i18n

import (
	"encoding/json"
	"fmt"
	"os"
)

const DefaultLang = "en-US"

type Catalog struct {
	Messages map[string]map[string]string
}

func LoadCatalog(path string) (*Catalog, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read i18n catalog: %w", err)
	}

	messages := map[string]map[string]string{}
	if err := json.Unmarshal(content, &messages); err != nil {
		return nil, fmt.Errorf("parse i18n catalog: %w", err)
	}

	if _, ok := messages[DefaultLang]; !ok {
		return nil, fmt.Errorf("i18n catalog missing default language %q", DefaultLang)
	}

	return &Catalog{Messages: messages}, nil
}

func (c *Catalog) Translate(lang, key string) string {
	if langMap, ok := c.Messages[lang]; ok {
		if message, ok := langMap[key]; ok {
			return message
		}
	}

	if defaultMap, ok := c.Messages[DefaultLang]; ok {
		if message, ok := defaultMap[key]; ok {
			return message
		}
	}

	return key
}
