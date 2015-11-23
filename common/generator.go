package common

type Generator interface {
	Generate() (uint64, error)
}