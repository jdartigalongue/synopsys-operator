package utils

// GetVersionLabel will return the label including the version
func  GetVersionLabel(name string, version string) map[string]string {
	m := GetLabel(name)
	m["version"] = version
	return m
}

// GetLabel will return the label
func  GetLabel(name string) map[string]string {
	return map[string]string{
		"app":       "blackduck",
		"component": name,
	}
}