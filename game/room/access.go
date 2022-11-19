package room

//go:generate go run golang.org/x/tools/cmd/stringer -type=Access

// Access represents a Room's access type (e.g. open, closed, password protected).
type Access int

const (
	Open Access = iota
	Closed
	Password
)

// AccessTypes returns a slice of all the possible Access values.
func AccessTypes() []Access {
	return []Access{
		Open,
		Closed,
		Password,
	}
}
