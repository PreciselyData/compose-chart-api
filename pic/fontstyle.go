package pic

// FontStyle represents a font style from Designer/Generate.
type FontStyle struct {
	*FontResource
	Color
	Underline bool
}
