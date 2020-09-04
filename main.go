package main

import (
    "goendpoint/controllers"
    "net/http"
    "log"
    "strconv"
)

func main() {
    res := controllers.HandleConsoleInput()

    controllers.AttachHandlers(res.FileName)

    log.Println("Server started on port " + strconv.Itoa(res.Port))
    http.ListenAndServe(":" + strconv.Itoa(res.Port), nil)
}
