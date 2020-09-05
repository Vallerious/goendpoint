package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"goendpoint/services"
	"io/ioutil"
	"net/http"
	"strings"
	"encoding/base64"
)

var AuthUser = ""
var AuthSecret = ""

func AttachHandlers(resource string) {
	getAllHandler := func(resp http.ResponseWriter, req *http.Request) (b []byte) {
		logs.Info("Hmmm...you want all the records? Do you know what you are doing?")
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
		logs.Info("Ohhhh...so you want to add something, right? Alrighty...")
		return upsertHandler(resp, req, func(id string, resource string, incomingData map[string]interface{}) (r []byte, fe error) {
			return services.Add(resource, incomingData)
		})
	}

	updateResource := func(resp http.ResponseWriter, req *http.Request) (b []byte) {
		logs.Info("Ohhhh...so you want to update something, right? Alrighty...")
		return upsertHandler(resp, req, func(id string, resource string, incomingData map[string]interface{}) (r []byte, fe error) {
			return services.Update(req.Header.Get("id"), resource, incomingData)
		})
	}

	dispatcher := func(resp http.ResponseWriter, req *http.Request) {
		var resData []byte

		logs.Info("The URL that you are calling is: " + req.URL.Path)

		if AuthUser != "" {
			logs.Info("Authentication enabled with user: " + AuthUser + " and secret " + AuthSecret)

			authorizationHeader := req.Header.Get("Authorization")

			logs.Info("Who are you!? " + authorizationHeader)

			headerValueParts := strings.Split(authorizationHeader, " ")

			if len(headerValueParts) != 2 {
				http.Error(resp, "Invalid basic auth payload", http.StatusBadRequest)
				return
			}

			base64Part := headerValueParts[1]
			decodedBase64, authDecodeErr := base64.StdEncoding.DecodeString(base64Part)

			if authDecodeErr != nil {
				http.Error(resp, authDecodeErr.Error(), http.StatusBadRequest)
				return
			}

			logs.Info("Hmmm I know your credentials....They are: " + string(decodedBase64))

			userSecretArgs := strings.Split(string(decodedBase64), ":")

			if len(userSecretArgs) != 2 {
				http.Error(resp, "Invalid basic auth payload", http.StatusBadRequest)
				return
			}

			if userSecretArgs[0] != AuthUser || userSecretArgs[1] != AuthSecret {
				http.Error(resp, "Invalid credentials", http.StatusUnauthorized)
				return
			}
		}

		// Used for server's internal purposes to find the resource in update and delete operations.
		pathSplit := strings.Split(req.URL.Path, "/")
		if len(pathSplit) == 3 {
			logs.Info("Let's see if we have a resource id: " + pathSplit[2])
			req.Header.Set("id", pathSplit[2])
		}

		if pathSplit[1] != resource {
			http.Error(resp, "Page not found", http.StatusNotFound)
			return
		}

		logs.Info("The request method is: " + req.Method)
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

	http.HandleFunc("/", dispatcher)
}

