package lib 

import(
	"unicode"
	"strings"
)

type extractor struct {
	password string
}

// Extracts lower case, upper case, special chars & numeric chars and groups them into maps
func (e *extractor) GetStringGroups() map[string]string {
	ret := make(map[string]string)
	ret["uppercase"] = ""
	ret["lowercase"] = ""
	ret["special"] = ""
	ret["number"] = ""
	for _, character := range e.password {
		if unicode.IsLower(character) {
			ret["lowercase"] = ret["lowercase"] + string(character)
		} else if unicode.IsUpper(character) {
			ret["uppercase"] = ret["uppercase"] + string(character)
		} else if unicode.IsSymbol(character) || unicode.IsSpace(character) || unicode.IsPunct(character) {
			ret["special"] = ret["special"] + string(character)
		} else if unicode.IsNumber(character) {
			ret["number"] = ret["number"] + string(character)
		}

	}
	return ret
}

func (e *extractor) GenerateCombinations(stringGroups map[string]string) []string {
	ret := make([]string, 0)
	combinations := []string{
		"lowercase,number,special",
		"special,uppercase,number",
		"number,special",
		"special,uppercase,lowercase",
	}

	for _, value := range combinations {
		s := strings.Split(value, ",")
		strGroupValue := ""
		for _, groupType := range s {
			groupValue, exists := stringGroups[groupType]
			if exists {
				strGroupValue = strGroupValue + groupValue
			}
		}
		ret = append(ret, strGroupValue)
	}
	return ret
}

func Extractor(password string) extractor {
	extractor := extractor{
		password: password,
	}
	return extractor
}