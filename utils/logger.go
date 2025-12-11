package utils

import (
    "log"
    "os"
)

// Logger levels
const (
    LevelError = "ERROR"
    LevelWarn  = "WARN"
    LevelInfo  = "INFO"
    LevelDebug = "DEBUG"
)

// Logger provides structured logging
type Logger struct {
    level string
    log   *log.Logger
}

// NewLogger creates a new logger with the specified level
func NewLogger(level string) *Logger {
    return &Logger{
        level: level,
        log:   log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
    }
}

// Info logs info level messages
func (l *Logger) Info(msg string, keyvals ...interface{}) {
    if l.shouldLog(LevelInfo) {
        l.log.Printf("[INFO] %s %v", msg, keyvals)
    }
}

// Warn logs warning level messages
func (l *Logger) Warn(msg string, keyvals ...interface{}) {
    if l.shouldLog(LevelWarn) {
        l.log.Printf("[WARN] %s %v", msg, keyvals)
    }
}

// Error logs error level messages
func (l *Logger) Error(msg string, err error, keyvals ...interface{}) {
    if l.shouldLog(LevelError) {
        l.log.Printf("[ERROR] %s: %v %v", msg, err, keyvals)
    }
}

// Debug logs debug level messages
func (l *Logger) Debug(msg string, keyvals ...interface{}) {
    if l.shouldLog(LevelDebug) {
        l.log.Printf("[DEBUG] %s %v", msg, keyvals)
    }
}

// shouldLog determines if a message should be logged based on level
func (l *Logger) shouldLog(messageLevel string) bool {
    levels := map[string]int{
        LevelError: 1,
        LevelWarn:  2,
        LevelInfo:  3,
        LevelDebug: 4,
    }
    
    currentLevel, exists := levels[l.level]
    if !exists {
        currentLevel = 3 // Default to INFO
    }
    
    messageLogLevel, exists := levels[messageLevel]
    if !exists {
        return false
    }
    
    return messageLogLevel <= currentLevel
}
