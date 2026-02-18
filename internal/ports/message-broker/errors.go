package messagebroker

type InvalidConfiguration string

func (err InvalidConfiguration) Error() string {
	return string(err)
}
