package config

type Configuration struct {
	StockPort int16
	CostPort  int16
	MongoURL  string
}

var config *Configuration

func Get() *Configuration {
	if config == nil {
		config = &Configuration{
			StockPort: 3344,
			CostPort:  3345,
			MongoURL:  "mongodb://localhost:27017",
		}
	}

	return config
}
