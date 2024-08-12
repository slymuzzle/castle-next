package roomtype

type Type string

const (
	Personal Type = "personal"
	Group    Type = "group"
)

// Values provides list valid values for Enum.
func (Type) Values() (kinds []string) {
	for _, s := range []Type{Personal, Group} {
		kinds = append(kinds, string(s))
	}
	return
}
