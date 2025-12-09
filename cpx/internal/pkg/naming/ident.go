package naming

import (
	"strings"
	"unicode"
)

// SafeIdent converts an arbitrary project name into a valid C++ identifier.
func SafeIdent(name string) string {
	if name == "" {
		return "project"
	}
	var b strings.Builder
	for i, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			if i == 0 && unicode.IsDigit(r) {
				b.WriteByte('_')
			}
			b.WriteRune(r)
		} else {
			b.WriteByte('_')
		}
	}
	if b.Len() == 0 {
		return "project"
	}
	return b.String()
}

// SafeIdentUpper returns the uppercase version of SafeIdent.
func SafeIdentUpper(name string) string {
	return strings.ToUpper(SafeIdent(name))
}

// SafeIdentTitle returns SafeIdent with the first letter capitalized.
func SafeIdentTitle(name string) string {
	id := SafeIdent(name)
	if id == "" {
		return "Project"
	}
	return strings.ToUpper(id[:1]) + id[1:]
}

