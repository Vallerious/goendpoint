package services

import (
	"encoding/json"
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

func loadResourceFile(resource string) (s []byte, e error) {
	s, e = ioutil.ReadFile(filepath.Join("db", resource + ".json"))
	return
}
