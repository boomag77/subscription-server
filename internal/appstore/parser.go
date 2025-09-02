package appstore

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type jwsHeader struct {
	Alg string   `json:"alg"`
	Typ string   `json:"typ,omitempty"`
	X5c []string `json:"x5c,omitempty"`
}

type DecodedJWS struct {
	HeaderBytes    []byte
	PayloadBytes   []byte
	SignatureBytes []byte
}

type AppStoreNotification struct {
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

type ClientNotification struct {
	BundleID              string `json:"bundleId"`
	AppAccountToken       string `json:"appAccountToken,omitempty"`
	SignedTransactionInfo string `json:"signedTransactionInfo"`
}

/*
Transaction fields:
	ExpiresDateMS - 	ms since epoch The UNIX time, in milliseconds, an auto-renewable subscription purchase expires or renews.
	RevocationDateMS - 	ms since epoch The UNIX time, in milliseconds, that the App Store refunded the transaction or revoked
						it from Family Sharing.
	RevocationReason - 	The reason the transaction was revoked:
			0 The App Store refunded the transaction on behalf of the customer for other reasons, for example, an accidental purchase.
			1 The App Store refunded the transaction on behalf of the customer due to an actual or perceived issue within your app.
*/

type Transaction struct {
	OriginalTransactionID string `json:"originalTransactionId"`
	TransactionID         string `json:"transactionId"`
	ProductID             string `json:"productId"`
	ExpiresDateMS         *int64 `json:"expiresDate,omitempty"`
	RevocationDateMS      *int64 `json:"revocationDate,omitempty"`
	RevocationReason      *int   `json:"revocationReason,omitempty"`
}

/*
RenewalInfo fields:
	AutoRenewStatus
		0 Automatic renewal is off. The customer has turned off automatic renewal for the subscription, and it won’t renew at the end
		of the current subscription period.
		1 Automatic renewal is on. The subscription renews at the end of the current subscription period.
	ExpirationIntent
		1 The customer canceled their subscription.
		2 Billing error; for example, the customer’s payment information is no longer valid.
		3 The customer didn’t consent to an auto-renewable subscription price increase that requires their consent,
			or to a subscription offer conversion that requires their consent, so the subscription expired. For more information
			about subscription price increases that require customer consent, see Auto-renewable subscription price increase
			thresholds. For more information about offer conversions that require customer consent, see Consent for subscription
			offer conversions.
		4 The product wasn’t available for purchase at the time of renewal.
		5 The subscription expired for some other reason.
*/

type RenewalInfo struct {
	AutoRenewStatus          *int   `json:"autoRenewStatus,omitempty"`
	ExpirationIntent         *int   `json:"expirationIntent,omitempty"`
	IsInBillingRetryPeriod   *bool  `json:"isInBillingRetryPeriod,omitempty"`
	GracePeriodExpiresDateMS *int64 `json:"gracePeriodExpiresDate,omitempty"`
}

func parseTransaction(signedTransaction string) (*Transaction, error) {

	decodedTransaction, err := decodeSignedJWS(signedTransaction)
	if err != nil {
		return nil, fmt.Errorf("failed to decode transaction: %w", err)
	}
	var transaction Transaction
	if err := json.Unmarshal(decodedTransaction.PayloadBytes, &transaction); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction: %w", err)
	}
	return &transaction, nil
}

func parseRenewalInfo(signedRenewalInfo string) (*RenewalInfo, error) {

	decodedRenewalInfo, err := decodeSignedJWS(signedRenewalInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to decode renewal info: %w", err)
	}
	var renewalInfo RenewalInfo
	if err := json.Unmarshal(decodedRenewalInfo.PayloadBytes, &renewalInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal renewal info: %w", err)
	}
	return &renewalInfo, nil
}

func parseAppStoreNotification(body io.Reader) (*AppStoreNotification, error) {
	decodedSignedPayload, err := decodeRawNotification(body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode notification: %w", err)
	}

	decodedJWS, err := decodeSignedJWS(decodedSignedPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to decode signed JWS: %w", err)
	}

	var header jwsHeader
	if err := json.Unmarshal(decodedJWS.HeaderBytes, &header); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JWS header: %w", err)
	}

	var parsed AppStoreNotification
	if err := json.Unmarshal(decodedJWS.PayloadBytes, &parsed); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JWS payload: %w", err)
	}

	return &parsed, nil
}

func parseClientNotification(body io.Reader) (*ClientNotification, error) {

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	var clientNotification ClientNotification
	if err := json.Unmarshal(bodyBytes, &clientNotification); err != nil {
		return nil, fmt.Errorf("failed to unmarshal client notification: %w", err)
	}

	return &clientNotification, nil
}

func decodeSignedJWS(signed string) (*DecodedJWS, error) {
	if signed == "" {
		return nil, fmt.Errorf("signed payload is empty")
	}
	parts := strings.Split(signed, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWS format: want 3 parts")
	}

	// Validate the signature
	if err := validateSignedJWS(parts[0], parts[1], parts[2]); err != nil {
		return nil, fmt.Errorf("failed to validate JWS: %w", err)
	}

	// header raw bytes
	hdrBytes, err := decodeBase64String(parts[0])
	if err != nil {
		return nil, fmt.Errorf("decode header: %w", err)
	}
	// payload raw bytes
	plBytes, err := decodeBase64String(parts[1])
	if err != nil {
		return nil, fmt.Errorf("decode payload: %w", err)
	}
	// signature (raw bytes)
	sigBytes, err := decodeBase64String(parts[2])
	if err != nil {
		return nil, fmt.Errorf("decode signature: %w", err)
	}

	return &DecodedJWS{
		HeaderBytes:    hdrBytes,
		PayloadBytes:   plBytes,
		SignatureBytes: sigBytes,
	}, nil
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

func decodeBase64String(input string) ([]byte, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(input)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}
	return decoded, nil
}
