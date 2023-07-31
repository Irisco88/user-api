package httpserver

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
	"net/http"
	"os"
	"path/filepath"
)

func (uhs *UserHTTPServer) DownloadAvatarHandler(resp http.ResponseWriter, request *http.Request) {
	code := mux.Vars(request)["code"]
	fileName, userID, err := DecodeEncodedAvatar(code)
	if err != nil {
		uhs.log.Error("Failed to parse code",
			zap.Error(err),
			zap.String("code", code),
		)
		respondWithError(resp, http.StatusBadRequest, "Failed to parse code")
		return
	}
	tempFile, err := os.CreateTemp("", fmt.Sprintf("tempfile.*%s", filepath.Ext(fileName)))
	if err != nil {
		uhs.log.Error("Failed to create temp file", zap.Error(err))
		respondWithError(resp, http.StatusInternalServerError, "internal error")
		return
	}
	defer os.Remove(tempFile.Name())
	objectName := fmt.Sprintf("user%d/%s", userID, fileName)
	if e := uhs.minioClient.FGetObject(request.Context(), uhs.envConfig.MinioAvatarsBucket, objectName, tempFile.Name(), minio.GetObjectOptions{}); e != nil {
		if minioErr, ok := e.(minio.ErrorResponse); ok && minioErr.Code == "NoSuchKey" {
			respondWithError(resp, http.StatusNotFound, "not found")
			return
		}
		uhs.log.Error("failed to get object", zap.Error(e))
		respondWithError(resp, http.StatusInternalServerError, "internal error")
		return
	}
	uhs.log.Info("temp file created", zap.String("path", tempFile.Name()))
	http.ServeFile(resp, request, tempFile.Name())
}
