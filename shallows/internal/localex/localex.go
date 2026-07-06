// Package localex provides small helpers for working with locale/language strings.
package localex

import "golang.org/x/text/language"

// FirstDefined returns the first non-blank locale string, defaulting to
// the undetermined language tag ("und") when none are defined.
func FirstDefined(locales ...string) string {
	undefined := language.Und.String()
	for _, l := range locales {
		if l != "" && l != undefined {
			return l
		}
	}

	return undefined
}
