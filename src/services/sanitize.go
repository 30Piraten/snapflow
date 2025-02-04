package services

import "strings"

// Sanitize sanitizes a stringIt replaces spaces and
// other unsafe characters with underscore
func Sanitize(folderName string) string {
	// Replace spaces with underscore
	folderName = strings.ReplaceAll(folderName, " ", "_")

	// Remove or replace other unsafe characters (e.g., / \ : * ? " < > |)
	unsafeCharacters := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}

	for _, chars := range unsafeCharacters {
		folderName = strings.ReplaceAll(folderName, chars, "_")
	}

	// Convert to lower case
	folderName = strings.ToLower(folderName)

	return folderName
}
