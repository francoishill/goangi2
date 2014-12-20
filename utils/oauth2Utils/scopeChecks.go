package oauth2Utils

import (
	"strings"
)

func ScopeHasRequiredScope(actualScopesCSV, requiredScope string) bool {
	sanitizedRequiredScope := strings.Trim(requiredScope, " ")
	if strings.Contains(requiredScope, ",") {
		panic("Required scope cannot contain a comma character")
	}

	actualScopeList := strings.Split(strings.Trim(actualScopesCSV, " ,"), ",")
	for _, actualScope := range actualScopeList {
		if strings.EqualFold(actualScope, sanitizedRequiredScope) {
			return true
		}
	}

	return false
}
