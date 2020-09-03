package services

import (
	"goendpoint/models"
	"os"
	"errors"
	"io/ioutil"
	"path/filepath"
	"encoding/json"
)

func CreateSchema(schemaName string, js map[string]interface{}) error {
	_, err := os.Stat("db")

	if os.IsNotExist(err) {
		err := os.Mkdir("db", 0755)

		if err != nil {
			return errors.New("unable to create a folder to work with, check your permissions")
		}
	}

	newSchema := models.Schema{Headers: js, Data: make([]interface{}, 0, 0)}

	jsonStr, _ := json.Marshal(newSchema)

	e := ioutil.WriteFile(filepath.Join("db", schemaName), jsonStr, 0755)

	if e != nil {
		return errors.New("unable to create a file to work with, check your permissions")
	}

	return nil
}

func ValidateSchema(schemaName string, incomingData map[string]interface{}) error {
	schema, err := LoadSchema(schemaName)

	if err != nil {
		return err
	}

	for key, _ := range schema.Headers {
		if incomingData[key] == nil {
			return errors.New(key + " is a required property")
		}
	}

	return nil
}
