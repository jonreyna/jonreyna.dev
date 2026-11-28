package plate

// Issue is one problem [Validate] found.
type Issue struct {
	// Where names the item, such as "link L05" or "person justin".
	Where string
	// Message says what is wrong, in lowercase.
	Message string
}

// String returns the issue as "where: message".
func (i Issue) String() string { return i.Where + ": " + i.Message }
