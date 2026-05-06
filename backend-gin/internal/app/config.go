package app

import "os"

type Config struct {
	MySQLDSN        string
	CORSAllowOrigin string
	UploadDir       string
	PublicBaseURL   string
	AIServiceURL    string
}

func LoadConfig() Config {
	return Config{
		MySQLDSN:        getEnv("MYSQL_DSN", "root:root@tcp(mysql:3306)/selfcook?charset=utf8mb4&parseTime=True&loc=Local"),
		CORSAllowOrigin: getEnv("CORS_ALLOW_ORIGINS", "http://localhost:5173"),
		UploadDir:       getEnv("UPLOAD_DIR", "./uploads"),
		PublicBaseURL:   getEnv("PUBLIC_BASE_URL", ""),
		AIServiceURL:    getEnv("AI_SERVICE_URL", "http://ai-service:8001"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
