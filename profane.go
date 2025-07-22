package main

import (
	"slices"
	"strings"
)

var profanes = []string{
	"kerfuffle",
	"sharbert",
	"fornax",
}

const replacement = "****"

func replaceProfane(s string) string {
	words := strings.Split(s, " ")
	for i, word := range words {
		if slices.Contains(profanes, strings.ToLower(word)) {
			words[i] = replacement
		}
	}
	return strings.Join(words, " ")
}
