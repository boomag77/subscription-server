package appstore

import (
	// "encoding/base64"
	// "encoding/json"
	// "fmt"
	// "io"
)

// func decodeRawNotification(r io.Reader) (string, error) {
// 	var rawBody struct {
// 		SignedPayload string `json:"signedPayload"`
// 	}
// 	dec := json.NewDecoder(io.LimitReader(r, 1<<20)) // Limit to 1MB
// 	dec.DisallowUnknownFields()
// 	if err := dec.Decode(&rawBody); err != nil {
// 		return "", fmt.Errorf("failed to decode notification: %w", err)
// 	}
// 	if rawBody.SignedPayload == "" {
// 		return "", fmt.Errorf("missing signed payload")
// 	}
// 	return rawBody.SignedPayload, nil
// }


// func decodeBase64String(input string) ([]byte, error) {
// 	decoded, err := base64.RawURLEncoding.DecodeString(input)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to decode base64: %w", err)
// 	}
// 	return decoded, nil
// }