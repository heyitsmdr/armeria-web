package armeria

const (
	SettingBrief       string = "brief"
	SettingWrap               = "wrap"
	SettingMaxLines           = "lines"
	SettingScriptTheme        = "script_theme"
)

// ValidSettings returns all valid settings for a Character.
func ValidSettings() []string {
	return []string{
		SettingBrief,
		SettingWrap,
		SettingMaxLines,
		SettingScriptTheme,
	}
}

// SettingDesc is used to retrieve the description of a Character setting.
func SettingDesc(name string) string {
	switch name {
	case SettingBrief:
		return "Toggle short room descriptions when moving."
	case SettingWrap:
		return "Wrap room descriptions at this character length."
	case SettingMaxLines:
		return "Truncate main display after this many lines."
	case SettingScriptTheme:
		return "Theme to use for the mob script editor."
	}

	return ""
}

// SettingDefault is used as a fallback for setting values.
func SettingDefault(name string) string {
	switch name {
	case SettingBrief:
		return "false"
	case SettingWrap:
		return "80"
	case SettingMaxLines:
		return "100"
	case SettingScriptTheme:
		return "one_dark"
	}

	return ""
}

// SettingValidationString returns a validation string to use prior to storing a setting value.
func SettingValidationString(name string) string {
	switch name {
	case SettingWrap:
		return "num|min:40|max:200"
	case SettingMaxLines:
		return "num|min:50|max:500"
	case SettingScriptTheme:
		return "in:one_dark,gruvbox,nord_dark"
	}

	return ""
}

// SettingPermission returns a specific permission that is required to view/alter a specific setting.
func SettingPermission(name string) string {
	switch name {
	case SettingScriptTheme:
		return "CAN_BUILD"
	}

	return ""
}
