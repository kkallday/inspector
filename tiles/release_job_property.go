package tiles

import "strings"

type ReleaseJobProperty struct {
	Name    string
	Default interface{}
	Link    string
	Job     string
	Release string
}

type ReleaseJobProperties []ReleaseJobProperty

func (rjps ReleaseJobProperties) Find(name string) (ReleaseJobProperty, bool) {
	for _, property := range rjps {
		if strings.HasPrefix(name, property.Name) {
			return property, true
		}
	}

	return ReleaseJobProperty{}, false
}
