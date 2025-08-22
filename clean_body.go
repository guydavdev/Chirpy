package main

import (
	"slices"
	"strings"
)

func cleanBody(body string) string {
	profanes := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Fields(body)
	for i, word := range words {
		if slices.Contains(profanes, strings.ToLower(word)) {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
