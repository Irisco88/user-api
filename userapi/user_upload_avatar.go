package userapi

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

func (uhs *UserHTTPServer) UploadAvatarHandler() runtime.HandlerFunc {
	return func(resp http.ResponseWriter, request *http.Request, pathParams map[string]string) {
		request.ParseMultipartForm(10 << 20) // Set maximum form size (10 MB in this example)
		// Get user_id and file from the form data
		userID, err := strconv.ParseUint(request.FormValue("user_id"), 10, 32)
		if err != nil {
			respondWithError(resp, http.StatusBadRequest, "parse user_id failed")
			return
		}

		file, fileHeader, err := request.FormFile("file")
		if err != nil {
			uhs.log.Error("failed to get file", zap.Error(err))
			respondWithError(resp, http.StatusBadRequest, "Failed to retrieve file from form")
			return
		}
		defer file.Close()

		// Validate file size
		if fileHeader.Size > int64(uhs.envConfig.UserAvatarMaxSize) {
			respondWithError(resp, http.StatusBadRequest, "File size exceeds the maximum limit of 5 MB")
			return
		}

		// Validate file extension (only allow PNG and JPEG)
		fileExt := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if fileExt != ".png" && fileExt != ".jpeg" && fileExt != ".jpg" {
			respondWithError(resp, http.StatusBadRequest, "Invalid file extension. Only PNG and JPEG are allowed")
			return
		}

		// Generate a unique code for the uploaded picture
		uniqueCode := uuid.New().String()

		// Save the photo into MinIO with the unique code as the object name
		objectName := fmt.Sprintf("user%d/%s%s", userID, uniqueCode, fileExt)
		fileInfo, err := uhs.minioClient.PutObject(request.Context(),
			uhs.envConfig.MinioAvatarsBucket,
			objectName,
			file,
			fileHeader.Size,
			minio.PutObjectOptions{ContentType: fileHeader.Header.Get("Content-Type")})
		if err != nil {
			uhs.log.Info("Failed to save photo to MinIO", zap.Error(err))
			if minioErr, ok := err.(minio.ErrorResponse); ok {
				fmt.Println("MinIO Error Code:", minioErr.Code)
				fmt.Println("MinIO Error Message:", minioErr.Message)
			}
			respondWithError(resp, http.StatusInternalServerError, "internal error")
			return
		}
		uhs.log.Info("file uploaded",
			zap.String("key", fileInfo.Key),
			zap.String("checksum", fileInfo.ETag))
		// Respond with the unique code for the uploaded picture
		respondWithJSON(resp, http.StatusCreated, map[string]string{
			"checksum": fileInfo.ETag,
			"file":     fmt.Sprintf("%s%s", uniqueCode, fileExt),
		})
	}
}
