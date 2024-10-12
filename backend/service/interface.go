package service

import "net/http"

type Controller interface {
	UploadFile(w http.ResponseWriter, r *http.Request) (error, string)
	GetBaseInfo(w http.ResponseWriter, r *http.Request) (error, [][]float64)
	ExportCSVHandler(w http.ResponseWriter, r *http.Request) (error, [][]string)
}
