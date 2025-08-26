package appstore

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"subscription-server/internal/storage"
)

func HandleAppStoreNotification(w http.ResponseWriter, r *http.Request, storage storage.Storage) {
	// Decode the raw notification
	signedPayload, err := decodeRawNotification(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to decode notification: %v", err), http.StatusBadRequest)

		return
	}

	// Parse the signed payload
	parsed, err := parseSignedPayload(signedPayload)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse signed payload: %v", err), http.StatusBadRequest)
		return
	}

	// Process the notification
	if err := processNotification(parsed, storage); err != nil {
		http.Error(w, fmt.Sprintf("failed to process notification: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func processNotification(parsed *ParsedJWS, storage storage.Storage) error {
	// Implement your notification processing logic here
	return nil
}

func decodeRawNotification(r io.Reader) (string, error) {
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

type Transaction struct {
	OriginalTransactionID string `json:"originalTransactionId"`
	TransactionID         string `json:"transactionId"`
	ProductID             string `json:"productId"`
	// ms since epoch The UNIX time, in milliseconds, an auto-renewable subscription purchase expires or renews.
	ExpiresDateMS *int64 `json:"expiresDate,omitempty"`
	// ms since epoch The UNIX time, in milliseconds, that the App Store refunded the transaction or revoked it from Family Sharing.
	RevocationDateMS *int64 `json:"revocationDate,omitempty"`
	// 0 The App Store refunded the transaction on behalf of the customer for other reasons, for example, an accidental purchase.
	// 1 The App Store refunded the transaction on behalf of the customer due to an actual or perceived issue within your app.
	RevocationReason *int `json:"revocationReason,omitempty"` // 0/1
}

type RenewalInfo struct {
	AutoRenewStatus *int `json:"autoRenewStatus,omitempty"` // 0/1
	// 0 Automatic renewal is off. The customer has turned off automatic renewal for the subscription, and it won’t renew at the end
	// of the current subscription period.
	// 1 Automatic renewal is on. The subscription renews at the end of the current subscription period.
	ExpirationIntent *int `json:"expirationIntent,omitempty"` // 1..5
	// 1 The customer canceled their subscription.
	// 2 Billing error; for example, the customer’s payment information is no longer valid.
	// 3 The customer didn’t consent to an auto-renewable subscription price increase that requires their consent,
	// or to a subscription offer conversion that requires their consent, so the subscription expired. For more information
	// about subscription price increases that require customer consent, see Auto-renewable subscription price increase
	// thresholds. For more information about offer conversions that require customer consent, see Consent for subscription
	// offer conversions.
	// 4 The product wasn’t available for purchase at the time of renewal.
	// 5 The subscription expired for some other reason.
	IsInBillingRetryPeriod   *bool  `json:"isInBillingRetryPeriod,omitempty"`
	GracePeriodExpiresDateMS *int64 `json:"gracePeriodExpiresDate,omitempty"`
}

func parseSignedPayload(signed string) (*ParsedJWS, error) {
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
