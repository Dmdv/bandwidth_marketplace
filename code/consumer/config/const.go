package config

type Deployment byte

const (
	Development Deployment = iota
	TestNet     Deployment = iota
	MainNet     Deployment = iota
)

// Development returns true if Deployment == Development. Either returns false.
func (d Deployment) Development() bool {
	return d == Development
}
