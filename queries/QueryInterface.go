package queries

type Query interface {
	GetValue() int
	GetName() string
}
