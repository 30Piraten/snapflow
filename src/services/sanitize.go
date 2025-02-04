package services

import "strings"

// SanitizeFolderConstruct sanitizes a string to be used as a folder name
// in the S3. It replaces spaces and other unsafe characters with underscore
func SanitizeFolder(folderName string) string {
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
