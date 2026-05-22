package typed

func findPlaceholdersInText(text string) []string {
	return inTextPlaceholderRegexp.FindAllString(text, -1)
}
