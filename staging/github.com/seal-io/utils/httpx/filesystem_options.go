package httpx

type FSOption struct {
	listable bool
	embedded bool
}

func FSOptions() *FSOption {
	return &FSOption{}
}

// WithListable enables directory listing.
func (o *FSOption) WithListable() *FSOption {
	o.listable = true
	return o
}

// WithEmbedded enables embedded mode,
// which provides a fixed modified time.
func (o *FSOption) WithEmbedded() *FSOption {
	o.embedded = true
	return o
}
