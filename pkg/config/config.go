package config

import (
	"encoding/json"
	"os"
)

type serviceConfiguration struct {
	Port int16 `json:"port"`
}

type Configuration struct {
	Composition serviceConfiguration `json:"composition"`

	MongoURL        string `json:"mongoUrl"`
	MongoAuthSource string `json:"mongoAuthSource"`
	MongoUsername   string `json:"mongoUsername"`
	MongoPassword   string `json:"mongoPassword"`

	RabbitURL string `json:"rabbitUrl"`

	RedisURL      string `json:"redisUrl"`
	RedisPassword string `json:"redisPassword"`
	RedisDB       int    `json:"redisDb"`

	AuthEnabled bool   `json:"authEnabled"`
	AuthURL     string `json:"authUrl"`
}

var config *Configuration

func Get() Configuration {
	if config == nil {
		config = &Configuration{
			Composition: serviceConfiguration{
				Port: 3344,
			},

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

		file, err := os.Open("config.json")
		defer file.Close()
		if err == nil && file != nil {
			json.NewDecoder(file).Decode(config)
		}
	}

	return *config
}
