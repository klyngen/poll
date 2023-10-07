package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/klyngen/votomatic-3000/packages/backend/models"
)

var configurationError = errors.New("configuration error")

func wrapError(err error) error {
	if err == nil {
		return nil
	}
	return errors.Join(configurationError, err)
}

func ReadConfiguration(configurationName string) (*models.Configuration, error) {
	jsonFile, err := os.Open(configurationName)

	if err != nil {
		return nil, wrapError(err)
	}

	defer jsonFile.Close()

	jsonBytes, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return nil, wrapError(err)
	}

	var configuration models.Configuration

	err = json.Unmarshal(jsonBytes, &configuration)

	return &configuration, wrapError(err)
}
