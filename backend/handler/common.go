package handler

import (
	"abandonlgzj/logger"
	"abandonlgzj/service"
	"abandonlgzj/types"
	"encoding/csv"
	"encoding/json"
	"net/http"
)

var controller service.Controller = &service.TController{}

type SessionD struct {
	Session string `json:"session"`
}

const (
	maxRows = 10000
)

func Upload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	logger.Console.Info("upload db")
	resp := types.JSResp{
		Success: true,
	}
	err, id := controller.UploadFile(w, r)
	if err != nil {
		logger.Console.Error(err.Error())
		resp.Msg = err.Error()
		resp.Success = false
	}

	var ss SessionD
	ss.Session = id
	resp.Data = ss
	//w.WriteHeader(http.StatusCreated)
	//json.NewEncoder(w).Encode(resp)
	d, _ := json.Marshal(resp)
	w.Write(d)
}

func GetBaseInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	logger.Console.Info("get geo data")
	resp := types.JSResp{
		Success: true,
	}
	err, list := controller.GetBaseInfo(w, r)
	if err != nil {
		logger.Console.Error(err.Error())
		resp.Msg = err.Error()
		resp.Success = false
	}
	resp.Data = list
	d, _ := json.Marshal(resp)
	w.Write(d)
}

func ExportCSVHandler(w http.ResponseWriter, r *http.Request) {
	logger.Console.Info("export csv")
	err, results := controller.ExportCSVHandler(w, r)
	if err != nil {
		//logger.Console.Error(err.Error())
		resp := types.JSResp{
			Success: false,
			Msg:     err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
	// 根据结果数量决定返回的数据结构
	if len(results) < maxRows {
		// 数据量小于 1 万条，返回数据
		response := types.Response{
			Flag: "client",
			Data: results,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		// 数据量大于 1 万条，生成 CSV 内容并返回
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=footprint_data.csv")

		writer := csv.NewWriter(w)
		defer writer.Flush()

		// 写入 CSV 文件头
		header := []string{"geoTime", "latitude", "longitude", "altitude", "course",
			"horizontalAccuracy", "verticalAccuracy", "speed",
			"status", "activity", "network", "appStatus",
			"dayTime", "groupTime", "isSplit", "isMerge",
			"isAdd", "networkName"}
		if err := writer.Write(header); err != nil {
			http.Error(w, "Error writing header to CSV", http.StatusInternalServerError)
			return
		}

		// 写入查询结果到 CSV
		for _, record := range results {
			if err := writer.Write(record); err != nil {
				http.Error(w, "Error writing record to CSV", http.StatusInternalServerError)
				return
			}
		}
	}
}
