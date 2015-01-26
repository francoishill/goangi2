package translationUtils

import (
	"fmt"
	"strings"
)

type translationProvider struct {
	translationMaps []*translationsMap
}

func CreateTranslationProvider(translationMaps []*translationsMap) *translationProvider {
	return &translationProvider{
		translationMaps: translationMaps,
	}
}

func (this *translationProvider) formatExpandedMapAsValue(mapToFormat map[string]string) string {
	stringList := []string{}
	if mapToFormat != nil {
		for key, val := range mapToFormat {
			stringList = append(stringList, fmt.Sprintf(`"%s":"%s"`, key, val))
		}
	}
	return fmt.Sprintf(`{%s}`, strings.Join(stringList, ", "))
}

func (this *translationProvider) GetTranslation(dotJoinedKeys string, params map[string]string) string {
	for _, tmap := range this.translationMaps {
		if wasFound, result := tmap.GetTranslationWithArguments(dotJoinedKeys, params); wasFound {
			return result
		}
	}
	// In case unable to translate because the KEY was not found
	return fmt.Sprintf("%s%s", dotJoinedKeys, this.formatExpandedMapAsValue(params))
}
