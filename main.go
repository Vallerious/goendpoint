package main

import (
    "goendpoint/controllers"
    "net/http"
    "log"
    "strconv"
)

func main() {
    res := controllers.HandleConsoleInput()

    controllers.AttachHandlers(res.Response.Msg)

    log.Println("May I take your order, please? Please put it on the counter number: " + strconv.Itoa(res.Port))
    http.ListenAndServe(":" + strconv.Itoa(res.Port), nil)
}
