package validators

func TextValidator(text string) bool {
	switch {
	case
		lengthValidator(text):
		return true
	default:
		return false
	}
}

func lengthValidator(text string) bool {
	if len([]rune(text)) < 1573 {
		return true
	}

	return false
}
