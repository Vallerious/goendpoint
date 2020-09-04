package services

import (
	"encoding/json"
	"goendpoint/models"
	"io/ioutil"
	"path/filepath"
	"goendpoint/utils"
)

func GetAll(resource string) ([]byte, error) {
	data, err := LoadResourceFile(resource)

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

func Add(resource string, data map[string]interface{}) (r []byte, e error) {
	return updateDB(resource, data, func(s *models.Schema) {
		s.Data = append(s.Data, data)
	})
}

func Update(id string, resource string, data map[string]interface{}) (r []byte, e error) {
	return updateDB(resource, data, func(s *models.Schema) {
		var toBeUpdated *map[string]interface{}

		for _, entry := range s.Data {
			if entry["id"] == id {
				toBeUpdated = &entry
				break
			}
		}

		// For safety, let's retain the id to avoid tampering with it.
		toBeUpdatedId := (*toBeUpdated)["id"]

		*toBeUpdated = data
		(*toBeUpdated)["id"] = toBeUpdatedId
	})
}

func updateDB(resource string, data map[string]interface{}, f func(s *models.Schema)) (r []byte, e error) {
	resourceSchema, resourceSchemaLoadErr := LoadSchema(resource)

	if resourceSchemaLoadErr != nil {
		e = resourceSchemaLoadErr
		return
	}

	f(&resourceSchema)

	storeSchemaErr := storeSchema(resource, &resourceSchema)

	if storeSchemaErr != nil {
		e = storeSchemaErr
		return
	}

	incomingDataStr, _ := json.Marshal(data)

	r = incomingDataStr
	return
}

func LoadResourceFile(resource string) (s []byte, e error) {
	s, e = ioutil.ReadFile(filepath.Join("db", resource + ".json"))
	return
}

func persist(resource string, data []byte) (e error) {
	e = ioutil.WriteFile(filepath.Join("db", resource + ".json"), data, 0755)
	return
}

func LoadSchema(resource string) (s models.Schema, e error) {
	data, loadErr := LoadResourceFile(resource)

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
