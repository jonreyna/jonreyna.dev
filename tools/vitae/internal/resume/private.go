package resume

// Private is the git-ignored overlay, resume.private.yaml: fields that never appear in
// public output.
type Private struct {
	// Basics holds the private header fields.
	Basics BasicsPrivate `yaml:"basics"`
}

// ApplyPrivate merges the overlay at path into r. Fields the overlay leaves empty keep
// r's values.
func ApplyPrivate(r *Resume, path string) error {

	var p Private
	if err := decodeFile(path, &p); err != nil {
		return err
	}

	if p.Basics.Phone != "" {
		r.Basics.Phone = p.Basics.Phone
	}

	return nil
}
