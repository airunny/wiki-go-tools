package ielastic

import elastic "github.com/olivere/elastic/v7"

type Config struct {
	Source   string `json:"source"`
	Username string `json:"username"`
	Password string `json:"password"`
	Sniff    bool   `json:"sniff"`
	Debug    bool   `json:"debug"`
}

func NewElastic(c *Config) (*elastic.Client, error) {
	return elastic.NewClient(
		elastic.SetURL(c.Source),
		elastic.SetBasicAuth(c.Username, c.Password),
		elastic.SetSniff(c.Sniff),
		//elastic.SetInfoLog(NewLogWithLevel(LogLevelInfo)),
		elastic.SetErrorLog(NewLogWithLevel(LogLevelError)),
		//elastic.SetTraceLog(NewLogWithLevel(LogLevelTrace)),
		elastic.SetGzip(true),
	)
}
