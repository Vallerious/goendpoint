package main

import (
    "goendpoint/controllers"
    "net/http"
    "log"
)

func main() {
    res := controllers.HandleConsoleInput()

    controllers.AttachHandlers(res.Msg)

    log.Println("Server started...")
    http.ListenAndServe(":3000", nil)
}
