package applestore

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
)

type appleParser struct {
	decoder *appleDecoder
}

func NewAppleParser(d *appleDecoder) *appleParser {
	return &appleParser{
		decoder: d,
	}
}

type jwsHeader struct {
	Alg string   `json:"alg"`
	Typ string   `json:"typ,omitempty"`
	X5c []string `json:"x5c,omitempty"`
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

func (p *appleParser) ParseTransaction(signedTransaction string) (*Transaction, error) {

	txPayloadBytes, err := p.decoder.DecodeSignedJWS(signedTransaction)
	if err != nil {
		return nil, fmt.Errorf("failed to decode transaction: %w", err)
	}
	var transaction Transaction
	if err := json.Unmarshal(txPayloadBytes, &transaction); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction: %w", err)
	}
	return &transaction, nil
}

func (p *appleParser) ParseRenewalInfo(signedRenewalInfo string) (*RenewalInfo, error) {

	riPayloadBytes, err := p.decoder.DecodeSignedJWS(signedRenewalInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to decode renewal info: %w", err)
	}
	var renewalInfo RenewalInfo
	if err := json.Unmarshal(riPayloadBytes, &renewalInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal renewal info: %w", err)
	}
	return &renewalInfo, nil
}

func (p *appleParser) ParseAppStoreNotification(body io.Reader) (*AppStoreNotification, error) {
	decodedSignedPayload, err := p.decoder.DecodeRawNotification(body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode notification: %w", err)
	}

	payloadBytes, err := base64.StdEncoding.DecodeString(decodedSignedPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to decode signed JWS: %w", err)
	}

	var parsed AppStoreNotification
	if err := json.Unmarshal(payloadBytes, &parsed); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JWS payload: %w", err)
	}

	return &parsed, nil
}

func (p *appleParser) ParseClientNotification(body io.Reader) (*ClientNotification, error) {

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
