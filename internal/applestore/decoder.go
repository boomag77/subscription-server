package applestore

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"subscription-server/internal/contracts"
)

type appleDecoder struct {
	validator contracts.JWSValidator
}

func NewAppleDecoder(v contracts.JWSValidator) *appleDecoder {
	return &appleDecoder{
		validator: v,
	}
}

func (d *appleDecoder) DecodeSignedJWS(signed string) ([]byte, error) {
	if signed == "" {
		return nil, fmt.Errorf("signed payload is empty")
	}
	parts := strings.Split(signed, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWS format: want 3 parts")
	}

	// Validate the signature
	if err := d.validator.Validate(parts[0], parts[1], parts[2]); err != nil {
		return nil, fmt.Errorf("failed to validate JWS: %w", err)
	}

	// payload raw bytes
	plBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("decode payload: %w", err)
	}

	return plBytes, nil
}

func (d *appleDecoder) DecodeRawNotification(r io.Reader) (string, error) {
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
