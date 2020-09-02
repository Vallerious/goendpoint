package controllers

import (
	"flag"
	"goendpoint/utils"
	"goendpoint/models"
	"io/ioutil"
	"goendpoint/services"
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

	services.CreateSchema(fileFlag, jsonAsKeyValue)

	return models.Response{Status: 200, Msg: "Success!"}
}
