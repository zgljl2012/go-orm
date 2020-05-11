package fields

// Function Options Pattern

// FieldOptions options of field
type FieldOptions struct {
	PrimaryKey bool
	Length     int
	Null       bool
}

var defaultOptions = FieldOptions{
	PrimaryKey: false,
	Length:     100,
	Null:       true,
}

// FieldOption option setter
type FieldOption func(options *FieldOptions)

// WithPrimaryKey set a field be primary key
func WithPrimaryKey(set bool) FieldOption {
	return func(options *FieldOptions) {
		options.PrimaryKey = set
	}
}

// WithLength set the length
func WithLength(length int) FieldOption {
	return func(options *FieldOptions) {
		options.Length = length
	}
}

// WithNull set null
func WithNull(null bool) FieldOption {
	return func(options *FieldOptions) {
		options.Null = null
	}
}
