package config

type Configuration struct {
	Port     int16
	MongoURL string
}

var config *Configuration

func Get() *Configuration {
	if config == nil {
		config = &Configuration{
			Port:     3344,
			MongoURL: "mongodb://localhost:27017",
		}
	}

	return config
}
