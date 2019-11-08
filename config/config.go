package config

type Auth struct {
	Username string
	Password string
}

type Configuration struct {
	Port            int16
	MongoURL        string
	MongoAuthSource string
	MongoUsername   string
	MongoPassword   string
	RabbitURL       string
}

var config *Configuration

func Get() Configuration {
	if config == nil {
		config = &Configuration{
			Port:            3344,
			MongoURL:        "mongodb://localhost:27017",
			MongoAuthSource: "admin",
			MongoUsername:   "admin",
			MongoPassword:   "admin",
			RabbitURL:       "amqp://guest:guest@localhost:5672",
		}
	}

	return *config
}
