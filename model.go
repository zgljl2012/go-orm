package orm

// Those interfaces are the user can implement, some of these are required, e.g. ModelFields.

// ModelFields get all fields
type ModelFields interface {
	Fields() []Field // all fields
}
