package plate

// Check records that Jon read a link's passage and confirmed the card.
type Check struct {
	// On is the date of the check, as YYYY-MM-DD.
	On string `yaml:"on"`
	// Edition is the edition or translation the passage was checked in.
	Edition string `yaml:"edition"`
}
