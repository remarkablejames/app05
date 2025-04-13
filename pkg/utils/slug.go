package utils

import (
	"github.com/gosimple/slug"
	"strconv"
)

func GenerateUniqueSlug(title string, exists func(string) bool) string {
	baseSlug := slug.Make(title)
	finalSlug := baseSlug
	counter := 1

	for exists(finalSlug) {
		finalSlug = baseSlug + "-" + strconv.Itoa(counter)
		counter++
	}

	return finalSlug
}
