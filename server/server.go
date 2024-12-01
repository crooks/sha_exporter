package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

type GroupMetrics struct {
	Hash string `yaml:"hash"`
}

type FileMetrics struct {
	Path string `yaml:"path"`
	Hash string `yaml:"hash"`
}

// Metrics contains all the configuration settings.  Both Groups and Files maps
// are keyed by a string.  In the case of Groups, it should match the name of
// the group to be processed (E.g. wheel).  For Files, the key is a freeform
// shortname that relates to the filename in Path (E.g. sshdcfg for
// /etc/ssh/sshd_config)
type Metrics struct {
	Groups map[string]GroupMetrics `yaml:"groups"`
	Files  map[string]FileMetrics  `yaml:"files"`
}

func (met *Metrics) WriteMetricsYaml(filename string) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()
	data, err := yaml.Marshal(met)
	if err != nil {
		return
	}
	err = os.WriteFile(filename, data, 0644)
	return
}

func (met *Metrics) WriteMetricsJson(filename string) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()
	data, err := json.Marshal(met)
	if err != nil {
		return
	}
	err = os.WriteFile(filename, data, 0644)
	return
}

func (met *Metrics) ServApi(w http.ResponseWriter, req *http.Request) {
	data, err := json.Marshal(met)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s\n", string(data))
}

// newMetrics imports a yaml formatted config file into a Config struct
func NewMetrics(filename string) (*Metrics, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Metrics{}
	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}
