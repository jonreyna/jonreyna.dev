package resume

// Highlight is one résumé bullet.
type Highlight struct {
	// Lead is the optional bold lead-in, without its trailing period.
	Lead string `yaml:"lead"`
	// Text is the bullet itself, in inline markup.
	Text string `yaml:"text"`
	// Tags mark what the bullet is about; variants boost or exclude by tag. Each must
	// be declared in [Resume.Tags].
	Tags []string `yaml:"tags"`
	// Weight ranks the bullet from 1 to 100. When a variant trims, higher weights stay
	// longer.
	Weight int `yaml:"weight"`
}
