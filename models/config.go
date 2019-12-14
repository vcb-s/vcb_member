package models

type song struct {
	Name     string
	Duration int
}

// Songs config
type Songs struct {
	Song []song
}
