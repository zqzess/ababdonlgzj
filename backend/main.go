package main

import (
	"abandonlgzj/config"
	"abandonlgzj/console"
	"abandonlgzj/handler"
	"abandonlgzj/logger"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	commands = make(chan string)
	wg       sync.WaitGroup
)

func main() {
	logger.InitConsoleLogger(true, true, true, true, false)
	path, _ := os.Getwd()
	jsonFilePath := flag.String("config", filepath.Join(path, "abandonlgzj.json"), "Path to JSON configuration file")
	// 解析命令行参数, -config /path/to/your/file.json
	flag.Parse()
	s := config.GetInstance()
	s.Start(jsonFilePath)

	// 自动删除
	go autoRemove()

	// 启动html服务器
	go func() {
		frontPort := s.GetFrontPort()
		logger.Console.Info("html start on %d", frontPort)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", frontPort), nil); err != nil {
			logger.Console.Error("Error starting server:", err)
		}
	}()

	// 启动命令处理的 goroutine
	wg.Add(1)
	go console.GetInstance().CommandHandler(commands)

	// 启动控制台输入的 goroutine
	go console.GetInstance().ConsoleInput(commands)

	/**
	接口
	*/

	// 上传
	http.HandleFunc("/upload", handler.Upload)

	http.HandleFunc("/getData", handler.GetBaseInfo)

	http.HandleFunc("/export", handler.ExportCSVHandler)

	fs := http.FileServer(http.Dir("./dist"))
	http.Handle("/", fs)

	port := s.GetBackPort()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: corsMiddleware(http.DefaultServeMux),
	}

	logger.Console.Info("服务器正在%d端口上监听...", port)
	// 启动服务器
	err := server.ListenAndServe()
	if err != nil {
		logger.Console.Error(err.Error())
		return
	}
	// 等待命令处理完成
	wg.Wait()
}

func getHttpRemoteAddr(req *http.Request) string {
	addr := ""
	if req != nil {
		//优先从参数里获取远程IP地址
		addr = req.RemoteAddr
		if addr != "" {
			return addr
		}

		//从Http的Header里获取远程IP地址(通过Nginx网关会屏蔽真实的远程IP)
		addr = req.Header.Get("X-Forwarded-For")
		if addr != "" {
			return addr
		}

		//最后从Http连接获取远程IP地址
		addr = req.RemoteAddr
		if strings.Contains(addr, ":") {
			return strings.Split(addr, ":")[0]
		}
	}
	return addr
}

// 中间件函数，用于设置跨域请求的响应头
func corsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		addr := getHttpRemoteAddr(r)
		logger.Console.Warn("remote ip : %s", addr)
		// 允许的域名列表
		allowedOrigins := []string{
			"http://localhost:5175",
			"https://pic.zqzess.top.:10243",
			"https://pic.whitemoon.top:10244",
			"https://pic2.zqzess.top:10248",
			"https://pic2.whitemoon.top:10249",
			"http://127.0.0.1:5175",
		}

		// 检查请求的 Origin 头部，并设置相应的 Access-Control-Allow-Origin
		origin := r.Header.Get("Origin")
		logger.Console.Warn("client origin : %s", origin)
		var allowedOrigin string
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				allowedOrigin = allowed
				break
			}
		}

		// 如果请求的 Origin 不在允许的列表中，则不允许访问
		if allowedOrigin == "" {
			logger.Console.Warn("no origin")
			//http.Error(w, "CORS not allowed", http.StatusForbidden)
			//return
		}

		//w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		// 允许所有域名访问，可以根据实际需求进行修改
		w.Header().Set("Access-Control-Allow-Origin", "*") // 前端的地址 //指定允许访问该资源的外域 URI，对于携带身份凭证的请求不可使用通配符*
		// 允许的请求方法
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE") //指明实际请求所允许使用的 HTTP 方法
		// 允许的请求头字段
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, sessionID") //指明实际请求中允许携带的首部字段
		// 预检请求的缓存时间，单位为秒
		w.Header().Set("Access-Control-Max-Age", "3600")
		//logger.Console.Info("前端发送了option")
		// 如果是预检请求，直接返回
		if r.Method == http.MethodOptions {
			return
		}
		//logger.Console.Info("前端发送了post")
		// 调用下一个处理函数
		handler.ServeHTTP(w, r)
	})
}
