package utils

import (
	"notes/internal/models"
	"notes/pkg/validations"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func hasThreeConsecutiveRepeatedChars(s string) bool {
	runes := []rune(s)
	for i := 0; i < len(runes)-2; i++ {
		if runes[i] == runes[i+1] && runes[i] == runes[i+2] {
			return true
		}
	}
	return false
}
func ValidateAndFormatTitle(title string) (bool, string, error) {

	trimmedTitle := strings.TrimSpace(title)

	if trimmedTitle == "" {
		return false, "", validations.ErreEmptyTitle
	}
	if hasThreeConsecutiveRepeatedChars(trimmedTitle) {
		return false, "", validations.ErrRepeatedLetters
	}
	words := strings.Fields(trimmedTitle)
	for i, word := range words {
		words[i] = cases.Title(language.English).String(strings.ToLower(word))
	}
	formattedTitle := strings.Join(words, " ")
	if len(formattedTitle) > 50 || len(formattedTitle) < 5 {
		return false, "", validations.ErrCharactersExcess
	}

	return true, formattedTitle, nil
}

func ValidateAndFormatContent(content string) (bool, string, error) {
	trimmedContent := strings.TrimSpace(content)
	if trimmedContent == "" {
		return false, "", validations.ErrEmptyContent
	}
	trimmedContent = strings.ToLower(trimmedContent)
	if hasThreeConsecutiveRepeatedChars(trimmedContent) {
		return false, "", validations.ErrRepeatedLetters
	}

	if len(trimmedContent) > 70 || len(trimmedContent) < 10 {
		return false, "", validations.ErrCharactersContentExcess
	}
	return true, trimmedContent, nil
}

func ValidateAndFormatCategory(category string) (bool, string, error) {
	trimmedCategory := strings.TrimSpace(category)
	if trimmedCategory == "" {
		return false, "", validations.ErrEmptyCategory
	}
	if hasThreeConsecutiveRepeatedChars(trimmedCategory) {
		return false, "", validations.ErrRepeatedLetters
	}
	words := strings.Fields(trimmedCategory)
	for i, word := range words {
		words[i] = cases.Title(language.English).String(strings.ToLower(word))
	}
	formattedCategory := strings.Join(words, " ")
	if len(formattedCategory) > 30 || len(formattedCategory) < 2 {
		return false, "", validations.ErrCharactersExcessCat
	}
	return true, formattedCategory, nil
}

func CompareCategories(existingCats, updatedCats []models.Category) bool {
	if len(existingCats) != len(updatedCats) {
		return false
	}

	existingMap := make(map[string]bool)
	for _, cat := range existingCats {
		existingMap[cat.Name] = true
	}

	for _, cat := range updatedCats {
		if !existingMap[cat.Name] {
			return false
		}
	}

	return true
}

func ValidateAndFormatUsername(username string) (bool, *string, error) {
	if len(username) < 5 || len(username) > 20 || hasThreeConsecutiveRepeatedChars(username) {
		return false, nil, validations.ErrInvalidUser
	}
	formatted := strings.ToLower(username)
	return true, &formatted, nil
}
