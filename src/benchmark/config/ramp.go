package config

type Ramp interface {
	TargetRPS(duration int) int
}
