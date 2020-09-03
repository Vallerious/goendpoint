package controllers

import (
	"flag"
	"goendpoint/utils"
	"goendpoint/models"
	"io/ioutil"
	"goendpoint/services"
	"strings"
)

func HandleConsoleInput() models.Response {
	var fileFlag string

	flag.StringVar(&fileFlag, "f", "", "pass a file with json object of your model")
	flag.Parse()

	data, err := ioutil.ReadFile(fileFlag)

	if err != nil {
		return models.Response{Status: 400, Msg: err.Error()}
	}

	jsonAsKeyValue, err := utils.JsonToMap(data)

	if err != nil {
		return models.Response{Status: 400, Msg: err.Error()}
	}

	e := services.CreateSchema(fileFlag, jsonAsKeyValue)

	if e != nil {
		return models.Response{Status: 500, Msg: err.Error()}
	}

	resourceName := fileFlag[:strings.LastIndex(fileFlag, ".")]

	return models.Response{Status: 200, Msg: resourceName}
}
