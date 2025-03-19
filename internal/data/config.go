package data

type Config struct {
	Name               string
	ValidationFunction func(string) error
	Value              string
}
