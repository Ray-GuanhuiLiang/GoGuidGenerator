package common

type Generator interface {
	func Generate() int64
}