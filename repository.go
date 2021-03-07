// trainingsapp project main.go
package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Repository interface {
	Save() error
	Load() error
	String() string
}

type YAMLRepository struct {
	FileName string
	Data     interface{}
}

func NewYAMLRepository(fname string, p interface{}) *YAMLRepository {
	y := YAMLRepository{FileName: fname, Data: p}
	return &y
}

func (y YAMLRepository) String() string {
	data, err := yaml.Marshal(&y.Data)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return string(data)
}

func (y YAMLRepository) Save() error {
	data, err := yaml.Marshal(&y.Data)
	if err != nil {
		return fmt.Errorf("YAML file %v could not be marshalled: #%v ", y.FileName, err)
	}

	err = ioutil.WriteFile(y.FileName, data, 0777)
	if err != nil {
		return fmt.Errorf("YAML file %v could not be saved: #%v ", y.FileName, err)
	}
	return nil
}

func (y *YAMLRepository) Load() error {

	yamlFile, err := ioutil.ReadFile(y.FileName)
	if err != nil {
		return fmt.Errorf("YAML file %v could not be loaded: #%v ", y.FileName, err)
	}
	err = yaml.Unmarshal(yamlFile, y.Data)
	if err != nil {
		return fmt.Errorf("YAML file %v could not be unmarshalled: #%v ", y.FileName, err)
	}

	return nil
}
