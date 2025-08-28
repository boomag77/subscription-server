package appstore

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func validateSignedJWS(header string, payload string, signature string) error {
	if header == "" || payload == "" || signature == "" {
		return fmt.Errorf("header, payload, and signature must all be non-empty")
	}

	hdrBytes, err := decodeBase64String(header)
	if err != nil {
		return fmt.Errorf("decode header: %w", err)
	}
	var hdr jwsHeader
	if err := json.Unmarshal(hdrBytes, &hdr); err != nil {
		return fmt.Errorf("unmarshal header: %w", err)
	}
	if hdr.Alg != "ES256" {
		return fmt.Errorf("unsupported algorithm: %s", hdr.Alg)
	}
	if len(hdr.X5c) == 0 {
		return fmt.Errorf("missing X5c field")
	}
	leaf := hdr.X5c[0]

	der, err := base64.StdEncoding.DecodeString(leaf)
	if err != nil {
		return fmt.Errorf("decode leaf: %w", err)
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return fmt.Errorf("parse leaf cert: %w", err)
	}
	pubKey, ok := cert.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("not an ECDSA public key")
	}

	signingInput := fmt.Sprintf("%s.%s", header, payload)
	digest := sha256.Sum256([]byte(signingInput)) //hash the signing input

	sigBytes, err := decodeBase64String(signature)
	if err != nil {
		return fmt.Errorf("decode signature: %w", err)
	}

	if !ecdsa.VerifyASN1(pubKey, digest[:], sigBytes) {
		return fmt.Errorf("invalid JWS signature")
	}

	return nil
}
