package translationUtils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type translationsMap struct {
	topJsonObject interface{}
}

func CreateTranslationsMap(jsonString string) *translationsMap {
	var jsonObj interface{}
	err := json.Unmarshal([]byte(jsonString), &jsonObj)
	if err != nil {
		panic(err)
	}
	return &translationsMap{
		topJsonObject: jsonObj,
	}
}

func CreateTranslationsMap_FromJsonFile(filePath string) *translationsMap {
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return CreateTranslationsMap(string(fileBytes))
}

func (this *translationsMap) GetTranslationWithArguments(dotJoinedKeys string, params map[string]string) (wasFound bool, result string) {
	keysArray := strings.Split(dotJoinedKeys, ".")

	currentObj := this.topJsonObject

	finalIndex := len(keysArray) - 1
	for index, key := range keysArray {
		switch currentType := currentObj.(type) {
		case map[string]interface{}:
			tmpNestedObject, existsInMap := currentType[key]
			if !existsInMap {
				return false, ""
			}
			currentObj = tmpNestedObject

			if index == finalIndex {
				if objAsString, isAString := currentObj.(string); isAString {
					populatedResult := objAsString
					if params != nil {
						for key, val := range params {
							populatedResult = strings.Replace(populatedResult, fmt.Sprintf("{{%s}}", key), val, -1)
						}
					}
					return true, populatedResult
				} else {
					panic(fmt.Errorf("Invalid translation request, the value of key '%s' is a map instead of a string.", dotJoinedKeys))
				}
			}
			break
		default:
			panic(fmt.Errorf("Unmarshalled json has invalid type '%#T'.", currentObj))
		}
	}

	return false, ""
}
