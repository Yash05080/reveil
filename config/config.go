package config

import (
    "fmt"
    "os"
    "strconv"
    "github.com/joho/godotenv"
)

type Config struct {
    // Supabase Configuration
    SupabaseURL        string
    SupabaseAnonKey    string
    SupabaseServiceKey string
    
    // API Configuration
    Port     string
    LogLevel string
    
    // JWT Configuration
    JWTSecret string
    
    // Encryption Configuration
    MasterEncryptionKey string
    
    // Heavy Model Configuration (for later phases)
    HuggingFaceAPIKey   string
    HeavyModelTimeout   string
    
    // Queue Configuration
    QueueWorkers int
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
    // Load .env file if exists
    godotenv.Load()
    
    cfg := &Config{
        SupabaseURL:        getEnv("SUPABASE_URL", ""),
        SupabaseAnonKey:    getEnv("SUPABASE_ANON_KEY", ""),
        SupabaseServiceKey: getEnv("SUPABASE_SERVICE_KEY", ""),
        Port:               getEnv("API_PORT", "8080"),
        LogLevel:           getEnv("LOG_LEVEL", "info"),
        JWTSecret:          getEnv("JWT_SECRET", ""),
        MasterEncryptionKey: getEnv("MASTER_ENCRYPTION_KEY", ""),
        HuggingFaceAPIKey:  getEnv("HUGGINGFACE_API_KEY", ""),
        HeavyModelTimeout:  getEnv("HEAVY_MODEL_TIMEOUT", "30s"),
        QueueWorkers:       getEnvAsInt("QUEUE_WORKERS", 10),
    }
    
    // Validate required fields
    if cfg.SupabaseURL == "" {
        return nil, fmt.Errorf("SUPABASE_URL is required")
    }
    if cfg.SupabaseServiceKey == "" {
        return nil, fmt.Errorf("SUPABASE_SERVICE_KEY is required")
    }
    if cfg.JWTSecret == "" {
        return nil, fmt.Errorf("JWT_SECRET is required")
    }
    if cfg.MasterEncryptionKey == "" {
        return nil, fmt.Errorf("MASTER_ENCRYPTION_KEY is required")
    }
    
    return cfg, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intVal, err := strconv.Atoi(value); err == nil {
            return intVal
        }
    }
    return defaultValue
}
