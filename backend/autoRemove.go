package main

import (
	"abandonlgzj/logger"
	"abandonlgzj/service"
	"time"
)

func autoRemove() {
	for {
		logger.Console.Debug("自动检测会话是否过期...")
		time.Sleep(time.Minute) // 每分钟检查一次
		service.GetInstance().RemoveExpiredSessions()
	}
}
