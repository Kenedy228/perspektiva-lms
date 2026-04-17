package selectable

func mapItems(rawItems []ItemParams) []Item {
	items := make([]Item, 0, len(rawItems))

	for i := range rawItems {
		item := NewItem(rawItems[i])
		items = append(items, item)
	}

	return items
}
