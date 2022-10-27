package k8sbuilder

type WithOption string

const (
	Overwrite WithOption = "overwrite"
	OverwriteIfDefaultValue WithOption = "overwriteIfDefaultValue"
	Merge WithOption = "merge"
)


// IsOverwrite permit to know if i should overwrite or not, base on options
// Default to true
func IsOverwrite(opts []WithOption) bool {
	if len(opts) == 0 || opts[0] == Overwrite {
		return true
	}

	return false
}


// IsOverwriteIfDefaultValue permit to know if I need to overwrite only if not default value
// Default to false
func IsOverwriteIfDefaultValue(opts []WithOption) bool {
	if len(opts) > 0 && opts[0] == OverwriteIfDefaultValue {
		return true
	}

	return false
}

// IsMerge permit to know if I need to merge items.
// Default to false
func IsMerge(opts []WithOption) bool {
	if len(opts) > 0 && opts[0] == Merge {
		return true
	}

	return false
}

