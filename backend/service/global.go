package service

import (
	"abandonlgzj/logger"
	"fmt"
	"os"
	"sync"
	"time"
)

// Session 代表一个用户会话
type Session struct {
	ID         string
	File       string
	LastActive time.Time // 记录最后活动时间
}

// UploadedFileManager 单例管理上传文件列表
type UploadedFileManager struct {
	sessions map[string]*Session
	mu       sync.Mutex
	timeout  time.Duration // 会话过期时间
}

// 单例实例
var instance *UploadedFileManager
var once sync.Once

// GetInstance 获取 UploadedFileManager 的单例实例
func GetInstance() *UploadedFileManager {
	once.Do(func() {
		instance = &UploadedFileManager{
			sessions: make(map[string]*Session),
			timeout:  time.Duration(conf.GetRemoveTime()) * time.Minute,
		}
	})
	return instance
}

// AddFile 添加文件名到会话，替换原有文件
func (m *UploadedFileManager) AddFile(sessionID, fileName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 更新或创建会话
	session, exists := m.sessions[sessionID]
	if exists {
		// 删除旧文件
		if err := os.Remove(session.File); err != nil {
			return err // 删除旧文件失败
		}
		// 更新文件名和最后活动时间
		session.File = fileName
		session.LastActive = time.Now()
	} else {
		// 创建新的会话
		m.sessions[sessionID] = &Session{
			ID:         sessionID,
			File:       fileName,
			LastActive: time.Now(),
		}
	}
	return nil
}

// GetSession 根据 session ID 获取会话
func (m *UploadedFileManager) GetSession(sessionID string) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.sessions[sessionID]
}

// GetFileBySessionID 根据 session ID 获取对应的文件路径，并更新活动时间
func (m *UploadedFileManager) GetFileBySessionID(sessionID string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, exists := m.sessions[sessionID]
	if !exists {
		return "", fmt.Errorf("数据文件不存在，请先上传") // 返回文件不存在的错误
	}

	// 更新最后活动时间
	session.LastActive = time.Now()
	return session.File, nil
}

// RemoveExpiredSessions 移除过期的会话
func (m *UploadedFileManager) RemoveExpiredSessions() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id, session := range m.sessions {
		if time.Since(session.LastActive) > m.timeout {
			logger.Console.Info("会话已经过期，删除会话: %s", id)
			// 删除文件
			os.Remove(session.File)
			// 移除会话
			delete(m.sessions, id)
		}
	}
}
