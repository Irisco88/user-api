package httpserver

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/openfms/authutil"
	commonpb "github.com/openfms/protos/gen/common/v1"
	"go.uber.org/zap"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

func (uhs *UserHTTPServer) UploadAvatarHandler(resp http.ResponseWriter, request *http.Request) {
	claims, found := authutil.TokenClaimsFromCtx(request.Context())
	if !found {
		http.Error(resp, "get claims failed", http.StatusUnauthorized)
		return
	}
	// Get user_id and file from the form data
	userID, err := strconv.ParseUint(request.FormValue("user_id"), 10, 32)
	if err != nil {
		http.Error(resp, "parse user_id failed", http.StatusBadRequest)
		return
	}
	if !(claims.Role == commonpb.UserRole_USER_ROLE_ADMIN ||
		(claims.Role == commonpb.UserRole_USER_ROLE_NORMAL && claims.UserID == uint32(userID))) {
		http.Error(resp, "invalid access", http.StatusUnauthorized)
		return
	}

	file, fileHeader, err := request.FormFile("file")
	if err != nil {
		uhs.log.Error("failed to get file", zap.Error(err))
		http.Error(resp, "Failed to retrieve file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file size
	if fileHeader.Size > int64(uhs.envConfig.UserAvatarMaxSize) {
		http.Error(resp, "File size exceeds the maximum limit of 5 MB", http.StatusBadRequest)
		return
	}

	// Validate file extension (only allow PNG and JPEG)
	fileExt := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if fileExt != ".png" && fileExt != ".jpeg" && fileExt != ".jpg" {
		http.Error(resp, "Invalid file extension. Only PNG and JPEG are allowed", http.StatusBadRequest)
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
		uhs.log.Info("Failed to s ave photo to MinIO", zap.Error(err))
		http.Error(resp, "internal error", http.StatusInternalServerError)
		return
	}
	uhs.log.Info("file uploaded",
		zap.String("key", fileInfo.Key),
		zap.String("checksum", fileInfo.ChecksumSHA256))
	// Respond with the unique code for the uploaded picture
	respondWithJSON(resp, http.StatusCreated, map[string]string{
		"checksum": fileInfo.ETag,
		"file":     fmt.Sprintf("%s%s", uniqueCode, fileExt),
	})
}
