package property

// NormalizeKey maps user-facing property names to internal keys.
func NormalizeKey(key string) string {
	switch key {
	case "color":
		return "font-color"
	case "bg":
		return "shading"
	default:
		return key
	}
}
