package controllers

import (
	"encoding/json"
	"goendpoint/services"
	"net/http"
	"io/ioutil"
)

func AttachHandlers(resource string) {
	getAllHandler := func(resp http.ResponseWriter, req *http.Request) (b []byte) {
		allItems, err := services.GetAll(resource)

		if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}

		return allItems
	}

	upsertHandler := func(resp http.ResponseWriter, req *http.Request, f func(id string, resource string, incomingData map[string]interface{}) (r []byte, fe error)) (b []byte) {
		// TODO: Replace with MaxBytesReader
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			http.Error(resp, err.Error(), http.StatusBadRequest)
			return
		}

		var incomingData map[string]interface{}
		unmarshallErr := json.Unmarshal(body, &incomingData)

		if unmarshallErr != nil {
			http.Error(resp, unmarshallErr.Error(), http.StatusBadRequest)
			return
		}

		validationErr := services.ValidateSchema(resource, incomingData)

		if validationErr != nil {
			http.Error(resp, validationErr.Error(), http.StatusBadRequest)
			return
		}

		r, saveToDiskErr := f("", resource, incomingData)

		if saveToDiskErr != nil {
			http.Error(resp, saveToDiskErr.Error(), http.StatusInternalServerError)
			return
		}

		return r
	}

	addResource := func(resp http.ResponseWriter, req *http.Request) (b []byte) {
		return upsertHandler(resp, req, func(id string, resource string, incomingData map[string]interface{}) (r []byte, fe error) {
			return services.Add(resource, incomingData)
		})
	}

	updateResource := func(resp http.ResponseWriter, req *http.Request) (b []byte) {
		return upsertHandler(resp, req, func(id string, resource string, incomingData map[string]interface{}) (r []byte, fe error) {
			return services.Update(resource, incomingData)
		})
	}

	dispatcher := func(resp http.ResponseWriter, req *http.Request) {
		var resData []byte

		switch req.Method {
		case http.MethodGet:
			resData = getAllHandler(resp, req)
		case http.MethodPost:
			resData = addResource(resp, req)
		case http.MethodPut:
			resData = updateResource(resp, req)
		}

		// Place to attach common headers
		resp.Header().Set("Content-Type", "application/json")
		resp.Write(resData)
	}

	http.HandleFunc("/" + resource, dispatcher)
}

