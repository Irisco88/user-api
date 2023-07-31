package httpserver

import (
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/openfms/authutil"
	commonpb "github.com/openfms/protos/gen/common/v1"
	"go.uber.org/zap"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func (uhs *UserHTTPServer) DownloadAvatarHandler(resp http.ResponseWriter, request *http.Request) {
	claims, found := authutil.TokenClaimsFromCtx(request.Context())
	if !found {
		http.Error(resp, "get claims failed", http.StatusUnauthorized)
		return
	}
	queryParams := request.URL.Query()
	userIDStr := queryParams.Get("user_id")
	fileName := queryParams.Get("file")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		// Handle the error if the conversion fails
		http.Error(resp, "Invalid user_id parameter", http.StatusBadRequest)
		return
	}
	uhs.log.Info("new request", zap.String("file", fileName), zap.Uint64("userID", userID))
	if !(claims.Role == commonpb.UserRole_USER_ROLE_ADMIN ||
		(claims.Role == commonpb.UserRole_USER_ROLE_NORMAL && claims.UserID == uint32(userID))) {
		http.Error(resp, "invalid access", http.StatusUnauthorized)
		return
	}
	tempFile, err := os.CreateTemp("", fmt.Sprintf("tempfile.*%s", filepath.Ext(fileName)))
	if err != nil {
		uhs.log.Error("Failed to create temp file", zap.Error(err))
		http.Error(resp, "internal error", http.StatusInternalServerError)
	}
	//defer os.Remove(tempFile.Name())
	objectName := fmt.Sprintf("user%d/%s", userID, fileName)
	if e := uhs.minioClient.FGetObject(request.Context(), uhs.envConfig.MinioAvatarsBucket, objectName, tempFile.Name(), minio.GetObjectOptions{}); e != nil {
		if minioErr, ok := e.(minio.ErrorResponse); ok && minioErr.Code == "NoSuchKey" {
			http.Error(resp, "object not found", http.StatusNotFound)
			return
		}
		uhs.log.Error("failed to get object", zap.Error(e))
		http.Error(resp, "internal error", http.StatusInternalServerError)
		return
	}
	uhs.log.Info("temp file created", zap.String("path", tempFile.Name()))
	http.ServeFile(resp, request, tempFile.Name())
}
