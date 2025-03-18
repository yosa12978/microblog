package pkg

import (
	"fmt"
	"strings"
)

type ValidationError map[string]string

func (e ValidationError) Error() string {
	var s strings.Builder
	for key, value := range e {
		s.WriteString(fmt.Sprintf("%s: %s\n", key, value))
	}
	return s.String()
}
