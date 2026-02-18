package messagebroker

type InvalidConfigurationError string

func (e InvalidConfigurationError) Error() string {
	return string(e)
}
