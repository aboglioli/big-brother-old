package config

type Configuration struct {
	Port int16

	MongoURL        string
	MongoAuthSource string
	MongoUsername   string
	MongoPassword   string

	RabbitURL string

	RedisURL      string
	RedisPassword string
	RedisDB       int

	AuthEnabled bool
	AuthURL     string
}

var config *Configuration

func Get() Configuration {
	if config == nil {
		config = &Configuration{
			Port: 3344,

			MongoURL:        "mongodb://localhost:27017",
			MongoAuthSource: "admin",
			MongoUsername:   "admin",
			MongoPassword:   "admin",

			RabbitURL: "amqp://guest:guest@localhost:5672",

			RedisURL:      "localhost:6379",
			RedisPassword: "",
			RedisDB:       0,

			AuthEnabled: false,
			AuthURL:     "http://localhost:3000/v1/users/current",
		}
	}

	return *config
}