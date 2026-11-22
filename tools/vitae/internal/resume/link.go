package resume

// Link is a profile or site link in the header.
type Link struct {
	// Kind names the link's service, such as "linkedin", "github" or "site".
	Kind string `json:"kind" yaml:"kind"`
	// URL is the full https URL.
	URL string `json:"url" yaml:"url"`
	// Text is what prints, usually the URL without its scheme.
	Text string `json:"text" yaml:"text"`
}
