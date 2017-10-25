package misc

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v1"
)

type Log struct {
	Level       string
	Filename    string
	MaxSizeMB   int
	MaxBackups  int
	MaxAgeDays  int
	WriteStdout bool
	Json        bool
}
type Http struct {
	Address      string
	ReadTimeout  int
	WriteTimeout int
}

type Conf struct {
	RefreshInterval int
	StatsdHost      string
	Log             Log
	CoreURL         string
	PlayerURL       string
	Http            Http
	Token           string
}

func LoadConf(filename string) (*Conf, error) {
	cnf := Conf{}

	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(source, &cnf); err != nil {
		return nil, err
	}

	log.WithField("filename", filename).Info("Loaded configuration")

	return &cnf, nil
}
