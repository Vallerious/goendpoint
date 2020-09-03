package services

import (
	"encoding/json"
	"goendpoint/models"
	"io/ioutil"
	"path/filepath"
	"goendpoint/utils"
)

func GetAll(resource string) ([]byte, error) {
	data, err := loadResourceFile(resource)

	if err != nil {
		return []byte(""), err
	}

	jsonMap, e := utils.JsonToMap(data)

	if e != nil {
		return []byte(""), e
	}

	jsonStr, _ := json.Marshal(jsonMap["Data"])

	return jsonStr, nil
}

func Add(resource string, data []byte) (r []byte, e error) {
	var incomingData interface{}
	err := json.Unmarshal(data, &incomingData)

	if err != nil {
		e = err
		return
	}

	resourceSchema, resourceSchemaLoadErr := loadSchema(resource)

	if resourceSchemaLoadErr != nil {
		e = resourceSchemaLoadErr
		return
	}

	resourceSchema.Data = append(resourceSchema.Data, incomingData)

	storeSchemaErr := storeSchema(resource, &resourceSchema)

	if storeSchemaErr != nil {
		e = storeSchemaErr
		return
	}

	incomingDataStr, _ := json.Marshal(incomingData)

	r = incomingDataStr
	return
}

func loadResourceFile(resource string) (s []byte, e error) {
	s, e = ioutil.ReadFile(filepath.Join("db", resource + ".json"))
	return
}

func persist(resource string, data []byte) (e error) {
	e = ioutil.WriteFile(filepath.Join("db", resource + ".json"), data, 0755)
	return
}

func loadSchema(resource string) (s models.Schema, e error) {
	data, loadErr := loadResourceFile(resource)

	if loadErr != nil {
		e = loadErr
		return
	}

	e = json.Unmarshal(data, &s)

	return
}

func storeSchema(resource string, s *models.Schema) (e error) {
	jsStr, marshallingErr := json.Marshal(s)

	if marshallingErr != nil {
		e = marshallingErr
		return
	}

	writeToFileErr := persist(resource, jsStr)

	if writeToFileErr != nil {
		e = writeToFileErr
		return
	}

	return
}
