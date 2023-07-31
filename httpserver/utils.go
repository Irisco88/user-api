package httpserver

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
)

func GetEncodedAvatar(fileName string, userID uint32) string {
	u := url.Values{
		"f": []string{fileName},
		"u": []string{fmt.Sprintf("%v", userID)},
	}.Encode()
	return base64.URLEncoding.EncodeToString([]byte(u))
}

func DecodeEncodedAvatar(encodedStr string) (fileName string, userID uint32, err error) {
	decoded, err := base64.URLEncoding.DecodeString(encodedStr)
	if err != nil {
		return "", 0, err
	}
	u, err := url.QueryUnescape(string(decoded))
	if err != nil {
		return "", 0, err
	}
	queryValues, err := url.ParseQuery(u)
	if err != nil {
		return "", 0, err
	}
	fileNames, ok := queryValues["f"]
	if !ok || len(fileNames) != 1 {
		return "", 0, fmt.Errorf("fileName not found or multiple values present")
	}
	fileName = fileNames[0]
	userIDs, ok := queryValues["u"]
	if !ok || len(userIDs) != 1 {
		return "", 0, fmt.Errorf("userID not found or multiple values present")
	}
	userID64, err := strconv.ParseUint(userIDs[0], 10, 32)
	if err != nil {
		return "", 0, err
	}
	userID = uint32(userID64)
	return fileName, userID, nil
}
