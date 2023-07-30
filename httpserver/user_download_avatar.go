package httpserver

import (
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/openfms/authutil"
	commonpb "github.com/openfms/protos/gen/common/v1"
	"go.uber.org/zap"
	"io"
	"net/http"
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
	if !(claims.Role == commonpb.UserRole_USER_ROLE_ADMIN ||
		(claims.Role == commonpb.UserRole_USER_ROLE_NORMAL && claims.UserID == uint32(userID))) {
		http.Error(resp, "invalid access", http.StatusUnauthorized)
		return
	}
	objectName := fmt.Sprintf("user%d/%s", userID, fileName)
	object, err := uhs.minioClient.GetObject(request.Context(), uhs.envConfig.MinioAvatarsBucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		uhs.log.Error("failed to get object", zap.Error(err))
		http.Error(resp, "internal error", http.StatusInternalServerError)
		return
	}
	defer object.Close()
	// Retrieve the object's metadata to obtain the Content-Type
	objectInfo, err := object.Stat()
	if err != nil {
		if minioErr, ok := err.(minio.ErrorResponse); ok && minioErr.Code == "NoSuchKey" {
			http.Error(resp, "object not found", http.StatusNotFound)
			return
		}
		uhs.log.Error("Failed to retrieve object metadata", zap.Error(err))
		http.Error(resp, "internal error", http.StatusInternalServerError)
		return
	}
	// Set the appropriate Content-Type header for the response
	resp.Header().Set("Content-Type", objectInfo.ContentType)

	// Copy the object data to the response body
	_, err = io.Copy(resp, object)
	if err != nil {
		uhs.log.Error("Failed to copy object data to response", zap.Error(err))
		http.Error(resp, "internal error", http.StatusInternalServerError)
		return
	}
}
