package httpserver

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

func (uhs *UserHTTPServer) UploadAvatarHandler(resp http.ResponseWriter, request *http.Request) {
	request.ParseMultipartForm(10 << 20) // Set maximum form size (10 MB in this example)
	ctx := request.Context()
	// Get user_id and file from the form data
	userID, err := strconv.ParseUint(request.FormValue("user_id"), 10, 32)
	if err != nil {
		http.Error(resp, "parse user_id failed", http.StatusBadRequest)
		return
	}

	uhs.log.Info("new request", zap.Uint64("userID", userID))
	file, handler, err := request.FormFile("file")
	if err != nil {
		http.Error(resp, "Failed to retrieve file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file size
	if handler.Size > int64(uhs.envConfig.UserAvatarMaxSize) {
		http.Error(resp, "File size exceeds the maximum limit of 5 MB", http.StatusBadRequest)
		return
	}

	// Validate file extension (only allow PNG and JPEG)
	fileExt := strings.ToLower(filepath.Ext(handler.Filename))
	if fileExt != ".png" && fileExt != ".jpeg" && fileExt != ".jpg" {
		http.Error(resp, "Invalid file extension. Only PNG and JPEG are allowed", http.StatusBadRequest)
		return
	}

	// Generate a unique code for the uploaded picture
	uniqueCode := uuid.New().String()

	// Save the photo into MinIO with the unique code as the object name
	objectName := fmt.Sprintf("user123/%s%s", uniqueCode, fileExt)
	_ = objectName
	fileInfo, err := uhs.minioClient.PutObject(ctx,
		uhs.envConfig.MinioAvatarsBucket,
		objectName,
		file,
		handler.Size,
		minio.PutObjectOptions{ContentType: handler.Header.Get("Content-Type")})
	if err != nil {
		http.Error(resp, "Failed to s ave photo to MinIO", http.StatusInternalServerError)
		return
	}
	fmt.Printf("%#v\n", fileInfo)
	uhs.log.Info("file uploaded",
		zap.String("key", fileInfo.Key),
		zap.String("checksum", fileInfo.ChecksumSHA256))
	// Respond with the unique code for the uploaded picture
	respondWithJSON(resp, http.StatusCreated, map[string]string{
		"checksum": fileInfo.ETag,
		"code":     uniqueCode,
	})
}
