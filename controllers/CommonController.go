package controllers

import (
	"goendpoint/services"
	"net/http"
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

	http.HandleFunc("/" + resource, getAllHandler)
}

