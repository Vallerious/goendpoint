package controllers

import (
	"goendpoint/services"
	"net/http"
	"io/ioutil"
)

func AttachHandlers(resource string) {
	getAllHandler := func(resp http.ResponseWriter, req *http.Request) {
		allItems, err := services.GetAll(resource)

		if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}

		resp.Header().Set("Content-Type", "application/json")
		resp.Write(allItems)
	}

	addResource := func(resp http.ResponseWriter, req *http.Request) {
		// TODO: Replace with MaxBytesReader
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			http.Error(resp, err.Error(), http.StatusBadRequest)
			return
		}

		r, saveToDiskErr := services.Add(resource, body)

		if saveToDiskErr != nil {
			http.Error(resp, saveToDiskErr.Error(), http.StatusInternalServerError)
			return
		}

		resp.Header().Set("Content-Type", "application/json")
		resp.Write(r)
	}

	dispatcher := func(resp http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			getAllHandler(resp, req)
		case http.MethodPost:
			addResource(resp, req)
		}
	}

	http.HandleFunc("/" + resource, dispatcher)
}

