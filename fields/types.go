package fields

// Type file type
type Type int

const (
	// INT int
	INT Type = iota
	// FLOAT float
	FLOAT
	// CHAR char
	CHAR
	// BOOL boolean
	BOOL
	// DATETIME datetime
	DATETIME
	// UINT64 big int(Uint64)
	UINT64
)

func (t Type) String() string {
	switch t {
	case INT:
		return "INT"
	case FLOAT:
		return "FLOAT"
	case CHAR:
		return "CHAR"
	case BOOL:
		return "BOOL"
	case DATETIME:
		return "DATETIME"
	case UINT64:
		return "BIGINT"
	}
	return ""
}
