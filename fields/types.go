package fields

// Type file type
type Type int

const (
	// INT int
	INT Type = iota
	// CHAR char
	CHAR
	// BOOL boolean
	BOOL
	// DATETIME datetime
	DATETIME
)

func (t Type) String() string {
	switch t {
	case INT:
		return "INT"
	case CHAR:
		return "CHAR"
	case BOOL:
		return "BOOL"
	case DATETIME:
		return "DATETIME"
	}
	return ""
}
