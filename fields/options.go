package fields

// Function Options Pattern

// FieldOptions options of field
type FieldOptions struct {
	PrimaryKey bool
}

var defaultOptions = FieldOptions{
	PrimaryKey: false,
}

// FieldOption option setter
type FieldOption func(options *FieldOptions)

// WithPrimaryKey set a field be primary key
func WithPrimaryKey(set bool) FieldOption {
	return func(options *FieldOptions) {
		options.PrimaryKey = set
	}
}
