package utils

// Checks if a item exists in a slice
func Contains(slice []string, itemToSearch string) bool {
	for _, item := range slice {
		if item == itemToSearch {
			return true
		}
	}
	return false
}

// Checks if a key exists in a map
func KeyInMap(m map[string]struct{}, key string) bool {
	_, exists := m[key]
	return exists
}
