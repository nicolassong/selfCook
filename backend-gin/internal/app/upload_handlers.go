package app

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	gin "github.com/gin-gonic/gin"
)

type uploadImageResponse struct {
	URL      string `json:"url"`
	FileName string `json:"fileName"`
	Size     int64  `json:"size"`
}

func (s Server) publicBaseURL(c *gin.Context) string {
	if strings.TrimSpace(s.cfg.PublicBaseURL) != "" {
		return strings.TrimRight(strings.TrimSpace(s.cfg.PublicBaseURL), "/")
	}
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	if forwarded := c.GetHeader("X-Forwarded-Proto"); forwarded != "" {
		scheme = strings.TrimSpace(strings.Split(forwarded, ",")[0])
	}
	host := strings.TrimSpace(c.Request.Host)
	if host == "" {
		host = c.Request.Host
	}
	return fmt.Sprintf("%s://%s", scheme, host)
}

func (s Server) uploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		writeError(c, http.StatusBadRequest, 40001, "file is required", gin.H{"detail": err.Error()})
		return
	}

	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		writeError(c, http.StatusBadRequest, 40001, "only image files are allowed", nil)
		return
	}

	if file.Size > 8*1024*1024 {
		writeError(c, http.StatusBadRequest, 40001, "image size must be <= 8MB", nil)
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == "" {
		switch contentType {
		case "image/png":
			ext = ".png"
		case "image/webp":
			ext = ".webp"
		default:
			ext = ".jpg"
		}
	}

	allowed := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}
	if !allowed[ext] {
		writeError(c, http.StatusBadRequest, 40001, "unsupported image type", nil)
		return
	}

	dateDir := time.Now().Format("20060102")
	targetDir := filepath.Join(s.cfg.UploadDir, dateDir)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		writeError(c, http.StatusInternalServerError, 50000, "create upload dir failed", gin.H{"detail": err.Error()})
		return
	}

	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	targetPath := filepath.Join(targetDir, fileName)
	if err := c.SaveUploadedFile(file, targetPath); err != nil {
		writeError(c, http.StatusInternalServerError, 50000, "save image failed", gin.H{"detail": err.Error()})
		return
	}

	url := fmt.Sprintf("%s/uploads/%s/%s", s.publicBaseURL(c), dateDir, fileName)
	writeSuccess(c, uploadImageResponse{
		URL:      url,
		FileName: fileName,
		Size:     file.Size,
	})
}
