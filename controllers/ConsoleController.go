package controllers

import (
	"flag"
	"goendpoint/utils"
	"goendpoint/models"
	"io/ioutil"
	"goendpoint/services"
	"net/http"
	"strings"
)

func HandleConsoleInput() models.ConsoleArgsResponse {
	var fileFlag string
	var user string
	var secret string

	numbPtr := flag.Int("p", 3000, "listening port")
	flag.StringVar(&fileFlag, "f", "", "pass a file with json object of your model")
	flag.StringVar(&user, "u", "", "username used for basic auth (omit this to skip auth)")
	flag.StringVar(&secret, "s", "", "secret used for basic auth")
	flag.Parse()

	AuthUser = user
	AuthSecret = secret

	data, err := ioutil.ReadFile(fileFlag)

	if err != nil {
		return models.ConsoleArgsResponse{
			Response: models.Response{Status: http.StatusBadRequest, Msg: err.Error()},
		}
	}

	jsonAsKeyValue, err := utils.JsonToMap(data)

	if err != nil {
		return models.ConsoleArgsResponse{
			Response: models.Response{Status: http.StatusBadRequest, Msg: err.Error()},
		}
	}

	e := services.CreateSchema(fileFlag, jsonAsKeyValue)

	if e != nil {
		return models.ConsoleArgsResponse{
			Response: models.Response{Status: http.StatusInternalServerError, Msg: err.Error()},
		}
	}

	resourceName := fileFlag[:strings.LastIndex(fileFlag, ".")]

	return models.ConsoleArgsResponse{
		Response: models.Response{Status: 200, Msg: resourceName},
		FileName: fileFlag,
		Port: *numbPtr,
	}
}
