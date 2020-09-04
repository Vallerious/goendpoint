package models

type Response struct {
	Status int16
	Msg    string
}

type ConsoleArgsResponse struct {
	Response
	FileName string
	Port int
}
