package pic

type resolver interface {
	integer(s string) (int32, error)
	number(s string) (float64, error)
	numberFormat() NumberFormat
	fontResource(guid GUID) (*FontResource, error)
	fontStyle(guid GUID) (*FontStyle, error)
}
