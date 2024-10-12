package service

import (
	"abandonlgzj/config"
	"abandonlgzj/logger"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type TController struct {
}

var conf = config.GetInstance()

func (s *TController) UploadFile(w http.ResponseWriter, r *http.Request) (error, string) {
	w.Header().Set("Content-Type", "application/json")
	// 只允许 POST 请求
	if r.Method != http.MethodPost {
		return fmt.Errorf("不支持的方法"), ""
	}

	// 限制文件大小为 100 MB
	maxFileSize := conf.GetUploadSize() * 1024 * 1024 // 100 MB
	if r.ContentLength > maxFileSize {
		return fmt.Errorf("文件大小超过限制"), ""
	}

	fileContent, header, err := r.FormFile("file")
	if err != nil {
		return fmt.Errorf("读取文件失败: " + err.Error()), ""
	}
	defer fileContent.Close()

	// 检查文件类型
	fileName := header.Filename
	if !isValidSQLiteFile(fileName) {
		return fmt.Errorf("文件类型不正确，必须是 .db 或 .DB 文件"), ""
	}

	// 创建 tmpFile 目录（如果不存在）
	tmpDir := conf.GetTmpPath()
	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		return fmt.Errorf("创建目录失败: " + err.Error()), ""
	}

	// 获取会话 ID
	//sessionID := r.FormValue("sessionID")
	sessionID := r.Header.Get("sessionID")
	if sessionID == "" {
		sessionID = uuid.New().String() // 如未提供，则生成新的 session ID
	}

	// 生成唯一文件名，包含时间戳和 UUID
	//uniqueFileName := fmt.Sprintf("uploaded_%d_%s.db", time.Now().UnixNano(), uuid.New().String())

	// 拼接完整的文件路径
	filePath := filepath.Join(tmpDir, fmt.Sprintf("%s.DB", sessionID))
	//filePath := filepath.Join(tmpDir, uniqueFileName)

	// 检查是否存在该 session
	existingSession := GetInstance().GetSession(sessionID)
	if existingSession != nil {
		// 更新活动时间
		existingSession.LastActive = time.Now()
		filePath = existingSession.File // 使用已有的文件路径
	}

	out, err := os.Create(filePath)
	if err != nil {
		//http.Error(w, "Unable to create file", http.StatusInternalServerError)
		return fmt.Errorf("创建文件失败: " + err.Error()), ""
	}
	defer out.Close()

	if _, err := io.Copy(out, fileContent); err != nil {
		//http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return fmt.Errorf("保存文件失败: " + err.Error()), ""
	}

	// 将文件信息添加到全局列表中
	if err := GetInstance().AddFile(sessionID, filePath); err != nil {
		//http.Error(w, "文件更新失败: "+err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("文件更新失败: " + err.Error()), ""
	}

	w.WriteHeader(http.StatusOK)

	logger.Console.Info("文件上传成功，文件大小: %d 字节", header.Size)
	return nil, sessionID
}

func (s *TController) GetBaseInfo(w http.ResponseWriter, r *http.Request) (error, [][]float64) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		//http.Error(w, "不支持的方法", http.StatusMethodNotAllowed)
		return fmt.Errorf("不支持的方法"), nil
	}

	beginTime := r.URL.Query().Get("beginTime")
	endTime := r.URL.Query().Get("endTime")

	if beginTime == "" || endTime == "" {
		return fmt.Errorf("beginTime 和 endTime 不能为空"), nil
	}

	// 获取会话 ID
	sessionID := r.Header.Get("sessionID")
	logger.Console.Debug("session : %s", sessionID)
	if sessionID == "" {
		return fmt.Errorf("请先上传文件"), nil
	}

	file, err := GetInstance().GetFileBySessionID(sessionID)
	if err != nil {
		return err, nil
	}

	// 连接SQLite数据库
	// 注意：这里的`file::memory:?cache=shared`表示在内存中创建一个临时数据库，仅作为示例。
	// 实际使用中，你应该替换为实际的数据库文件路径。
	//db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return fmt.Errorf("server error"), nil
	}
	defer db.Close()

	// 查询数据
	query := "SELECT latitude, longitude FROM footprintContent3 WHERE datetime(geoTime / 1000, 'unixepoch') > ? AND datetime(geoTime / 1000, 'unixepoch') < ?;"
	rows, err := db.Query(query, beginTime, endTime)
	if err != nil {
		return err, nil
	}
	defer rows.Close()

	// 存储查询结果的二维切片
	var coordinates [][]float64

	for rows.Next() {
		var latitude float64
		var longitude float64
		if err := rows.Scan(&latitude, &longitude); err != nil {
			logger.Console.Error(err.Error())
		}
		if latitude == 0 && longitude == 0 {
			// 这条数据异常，geoTime = 978307200000
			continue
		}
		coordinates = append(coordinates, []float64{latitude, longitude})
	}

	if err := rows.Err(); err != nil {
		return err, nil
	}
	return nil, coordinates
}

func (s *TController) ExportCSVHandler(w http.ResponseWriter, r *http.Request) (error, [][]string) {
	if r.Method != http.MethodGet {
		//http.Error(w, "不支持的方法", http.StatusMethodNotAllowed)
		return fmt.Errorf("不支持的方法"), nil
	}

	beginTime := r.URL.Query().Get("beginTime")
	endTime := r.URL.Query().Get("endTime")

	if beginTime == "" || endTime == "" {
		return fmt.Errorf("beginTime 和 endTime 不能为空"), nil
	}

	// 获取会话 ID
	sessionID := r.Header.Get("sessionID") // 从请求中获取 session ID
	logger.Console.Debug("session : %s", sessionID)
	if sessionID == "" {
		return fmt.Errorf("请先上传文件"), nil
	}

	file, err := GetInstance().GetFileBySessionID(sessionID)
	if err != nil {
		logger.Console.Error(err.Error())
		return err, nil
	}

	// 设置数据库连接
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		logger.Console.Error(err.Error())
		return fmt.Errorf("server error"), nil
	}
	defer db.Close()

	// 执行 SQL 查询
	query := "SELECT geoTime, latitude, longitude, altitude, course, horizontalAccuracy, verticalAccuracy, speed, status, activity, network, appStatus, dayTime, groupTime, isSplit, isMerge, isAdd, networkName FROM footprintContent3 WHERE datetime(geoTime / 1000, 'unixepoch') > ? AND datetime(geoTime / 1000, 'unixepoch') < ?;"
	rows, err := db.Query(query, beginTime, endTime)
	if err != nil {
		logger.Console.Error(err.Error())
		return fmt.Errorf("server error"), nil
	}
	defer rows.Close()

	// 统计结果行数
	var results [][]string
	for rows.Next() {
		var geoTime float64
		var latitude, longitude, altitude, course, horizontalAccuracy, verticalAccuracy, speed float64
		var status, activity, network, appStatus, dayTime, groupTime string
		var isSplit, isMerge, isAdd bool
		var networkName sql.NullString

		if err := rows.Scan(&geoTime, &latitude, &longitude, &altitude, &course,
			&horizontalAccuracy, &verticalAccuracy, &speed, &status, &activity,
			&network, &appStatus, &dayTime, &groupTime, &isSplit, &isMerge,
			&isAdd, &networkName); err != nil {
			logger.Console.Error(err.Error())
			return fmt.Errorf("server error"), nil
		}

		// 将数据转换为字符串并存储
		record := []string{
			strconv.FormatFloat(geoTime, 'f', 6, 64),
			strconv.FormatFloat(latitude, 'f', 6, 64),
			strconv.FormatFloat(longitude, 'f', 6, 64),
			strconv.FormatFloat(altitude, 'f', 6, 64),
			strconv.FormatFloat(course, 'f', 6, 64),
			strconv.FormatFloat(horizontalAccuracy, 'f', 6, 64),
			strconv.FormatFloat(verticalAccuracy, 'f', 6, 64),
			strconv.FormatFloat(speed, 'f', 6, 64),
			status,
			activity,
			network,
			appStatus,
			dayTime,
			groupTime,
			strconv.FormatBool(isSplit),
			strconv.FormatBool(isMerge),
			strconv.FormatBool(isAdd),
			networkName.String,
		}
		//if !networkName.Valid {
		//	record[len(record)-1] = "" // 或者设置为其他默认值
		//}
		results = append(results, record)
	}

	// 检查 rows.Next() 的错误
	if err := rows.Err(); err != nil {
		logger.Console.Error(err.Error())
		return fmt.Errorf("server error"), nil
	}
	return nil, results
}

// isValidSQLiteFile 检查文件名是否以 .db 或 .DB 结尾
func isValidSQLiteFile(fileName string) bool {
	ext := filepath.Ext(fileName)
	return ext == ".db" || ext == ".DB"
}
