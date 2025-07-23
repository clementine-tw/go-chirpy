package main

import (
	"strings"
)

var profanes = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}

func replaceProfane(s string) string {
	const replacement = "****"

	words := strings.Split(s, " ")
	for i, word := range words {
		if _, ok := profanes[word]; ok {
			words[i] = replacement
		}
	}

	return strings.Join(words, " ")
}
