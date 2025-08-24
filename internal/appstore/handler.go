package appstore

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func DecodeRawNotification(r io.Reader) (string, error) {
	var rawBody struct {
		SignedPayload string `json:"signedPayload"`
	}
	dec := json.NewDecoder(io.LimitReader(r, 1<<20)) // Limit to 1MB
	dec.DisallowUnknownFields()
	if err := dec.Decode(&rawBody); err != nil {
		return "", fmt.Errorf("failed to decode notification: %w", err)
	}
	if rawBody.SignedPayload == "" {
		return "", fmt.Errorf("missing signed payload")
	}
	return rawBody.SignedPayload, nil
}

type jwsHeader struct {
	Alg string   `json:"alg"`
	Typ string   `json:"typ,omitempty"`
	X5c []string `json:"x5c,omitempty"`
}

type Notification struct {
	NotificationType string `json:"notificationType"`
	Subtype          string `json:"subtype,omitempty"`
	NotificationUUID string `json:"notificationUUID"`
	Version          string `json:"version"`
	SignedDate       int64  `json:"signedDate"`
	Data             struct {
		BundleID              string `json:"bundleId"`
		BundleVersion         string `json:"bundleVersion"`
		Environment           string `json:"environment"`
		AppAccountToken       string `json:"appAccountToken,omitempty"`
		SignedTransactionInfo string `json:"signedTransactionInfo"`
		SignedRenewalInfo     string `json:"signedRenewalInfo,omitempty"`
	} `json:"data"`
}

type ParsedJWS struct {
	Header    jwsHeader
	Payload   Notification
	Signature []byte
}

func ParseSignedPayload(signed string) (*ParsedJWS, error) {
	if signed == "" {
		return nil, fmt.Errorf("signed payload is empty")
	}
	parts := strings.Split(signed, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWS format: want 3 parts")
	}

	// header
	hdrBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("decode header: %w", err)
	}
	var hdr jwsHeader
	if err := json.Unmarshal(hdrBytes, &hdr); err != nil {
		return nil, fmt.Errorf("unmarshal header: %w", err)
	}

	// payload
	plBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("decode payload: %w", err)
	}
	var pl Notification
	if err := json.Unmarshal(plBytes, &pl); err != nil {
		return nil, fmt.Errorf("unmarshal payload: %w", err)
	}

	// signature (raw bytes)
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, fmt.Errorf("decode signature: %w", err)
	}

	return &ParsedJWS{Header: hdr, Payload: pl, Signature: sig}, nil
}
