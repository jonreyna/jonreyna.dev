package renderer

// DocumentHighlight is one bullet in a [DocumentRole].
type DocumentHighlight struct {
	// Lead is the bold lead-in, if any, without its trailing period.
	Lead string `json:"lead,omitempty"`
	// Text is the bullet.
	Text Text `json:"text"`
	// Tags are the bullet's tags, for consumers that filter further.
	Tags []string `json:"tags"`
}
