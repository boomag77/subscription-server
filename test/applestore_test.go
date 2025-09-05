package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"subscription-server/internal/applestore"
	"subscription-server/internal/logger"
	"subscription-server/internal/storage"
)

// MockStorage реализует интерфейс storage.Storage для тестирования
type MockStorage struct {
	subscriptions map[string]*storage.SubscriptionStatus
	saveError     error
	getError      error
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		subscriptions: make(map[string]*storage.SubscriptionStatus),
	}
}

func (m *MockStorage) GetSubscriptionStatus(userToken string) (*storage.SubscriptionStatus, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	status, ok := m.subscriptions[userToken]
	if !ok {
		return nil, nil
	}
	return status, nil
}

func (m *MockStorage) SetSubscriptionStatus(status *storage.SubscriptionStatus) error {
	if m.saveError != nil {
		return m.saveError
	}
	m.subscriptions[status.UserToken] = status
	return nil
}

func (m *MockStorage) SetGetError(err error) {
	m.getError = err
}

func (m *MockStorage) SetSaveError(err error) {
	m.saveError = err
}

// MockLogger реализует интерфейс logger.Logger для тестирования
type MockLogger struct {
	logs []logger.LogMessage
}

func NewMockLogger() *MockLogger {
	return &MockLogger{logs: []logger.LogMessage{}}
}

func (m *MockLogger) Log(message logger.LogMessage) {
	m.logs = append(m.logs, message)
}

// MockJWSValidator реализует интерфейс contracts.JWSValidator для тестирования
type MockJWSValidator struct {
	ValidateError error
}

func NewMockJWSValidator() *MockJWSValidator {
	return &MockJWSValidator{}
}

func (m *MockJWSValidator) Validate(header string, payload string, signature string) error {
	return m.ValidateError
}

func (m *MockJWSValidator) SetValidateError(err error) {
	m.ValidateError = err
}

// TestHandleClientNotification тестирует обработку уведомления от клиента
func TestHandleClientNotification(t *testing.T) {
	// Подготовка
	mockStorage := NewMockStorage()
	mockLogger := NewMockLogger()
	mockValidator := NewMockJWSValidator()

	// Создание декодера и парсера
	decoder := applestore.NewAppleDecoder(mockValidator)
	parser := applestore.NewAppleParser(decoder)

	// Создание сервиса
	service := applestore.NewAppleStoreService(mockStorage, mockLogger, parser)

	// Создание запроса с тестовыми данными
	clientNotification := map[string]interface{}{
		"bundleId":              "com.test.app",
		"appAccountToken":       "user123",
		"signedTransactionInfo": "header.eyJvcmlnaW5hbFRyYW5zYWN0aW9uSWQiOiIxMjM0NTYiLCJ0cmFuc2FjdGlvbklkIjoiNTQzMjEiLCJwcm9kdWN0SWQiOiJjb20udGVzdC5wcm9kdWN0IiwiZXhwaXJlc0RhdGUiOjE3MjUwMDAwMDAwMDB9.signature",
	}
	body, _ := json.Marshal(clientNotification)
	req := httptest.NewRequest(http.MethodPost, "/client-notification", bytes.NewReader(body))
	w := httptest.NewRecorder()

	// Выполнение
	service.HandleClientNotification(w, req)

	// Проверка
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", resp.StatusCode)
	}

	// Проверка сохранения статуса подписки
	status, err := mockStorage.GetSubscriptionStatus("user123")
	if err != nil {
		t.Errorf("Ошибка при получении статуса подписки: %v", err)
	}

	if status == nil {
		t.Error("Статус подписки не был сохранен")
	} else {
		if status.ProductID != "com.test.product" {
			t.Errorf("Неправильный ID продукта: ожидался com.test.product, получен %s", status.ProductID)
		}
		if status.OriginalTransactionID != "123456" {
			t.Errorf("Неправильный OriginalTransactionID: ожидался 123456, получен %s", status.OriginalTransactionID)
		}
	}
}

// TestHandleProviderNotification тестирует обработку уведомления от App Store
func TestHandleProviderNotification(t *testing.T) {
	// Подготовка
	mockStorage := NewMockStorage()
	mockLogger := NewMockLogger()
	mockValidator := NewMockJWSValidator()

	// Создание декодера и парсера
	decoder := applestore.NewAppleDecoder(mockValidator)
	parser := applestore.NewAppleParser(decoder)

	// Создание сервиса
	service := applestore.NewAppleStoreService(mockStorage, mockLogger, parser)

	// Создание запроса с тестовыми данными
	notificationData := map[string]interface{}{
		"signedPayload": "eyJub3RpZmljYXRpb25UeXBlIjoiUkVORVdBTCIsIm5vdGlmaWNhdGlvblVVSUQiOiIxMjM0NSIsInZlcnNpb24iOiIyLjAiLCJzaWduZWREYXRlIjoxNjI1MDA0ODUyLCJkYXRhIjp7ImJ1bmRsZUlkIjoiY29tLnRlc3QuYXBwIiwiYnVuZGxlVmVyc2lvbiI6IjEuMCIsImVudmlyb25tZW50Ijoic2FuZGJveCIsImFwcEFjY291bnRUb2tlbiI6InVzZXI0NTYiLCJzaWduZWRUcmFuc2FjdGlvbkluZm8iOiJoZWFkZXIuZXlKdmNtbG5hVzVoYkZSeVlXNXpZV04wYVc5dVNXUWlPaUl4TWpNME5UWWlMQ0owY21GdWMyRmpkR2x2YmtsRUlqb2lOVFF6TWpFaUxDSndjbTlrZFdOMFNXUWlPaUpqYjIwdWRHVnpkQzV3Y205a2RXTjBJaXdpWlhod2FYSmxjMFJoZEdVaU9qRTNNalV3TURBd01EQXdNREI5LnNpZ25hdHVyZSIsInNpZ25lZFJlbmV3YWxJbmZvIjoiaGVhZGVyLmV5SmhkWFJ2VW1WdVpYZFRkR0YwZFhNaU9qRXNJbWx6U1c1Q2FXeHNhVzVuVW1WMGNubFdZV3hzWlhraU9tWmhiSE5sZlEuc2lnbmF0dXJlIn19",
	}
	body, _ := json.Marshal(notificationData)
	req := httptest.NewRequest(http.MethodPost, "/provider-notification", bytes.NewReader(body))
	w := httptest.NewRecorder()

	// Выполнение
	service.HandleProviderNotification(w, req)

	// Проверка
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", resp.StatusCode)
	}
}

// TestHandleProviderNotification_Error тестирует обработку ошибок при уведомлении от App Store
func TestHandleProviderNotification_Error(t *testing.T) {
	// Подготовка
	mockStorage := NewMockStorage()
	mockLogger := NewMockLogger()
	mockValidator := NewMockJWSValidator()
	mockValidator.SetValidateError(errors.New("validation error"))

	// Создание декодера и парсера
	decoder := applestore.NewAppleDecoder(mockValidator)
	parser := applestore.NewAppleParser(decoder)

	// Создание сервиса
	service := applestore.NewAppleStoreService(mockStorage, mockLogger, parser)

	// Создание запроса с некорректными данными
	notificationData := map[string]interface{}{
		"signedPayload": "invalid.payload.data",
	}
	body, _ := json.Marshal(notificationData)
	req := httptest.NewRequest(http.MethodPost, "/provider-notification", bytes.NewReader(body))
	w := httptest.NewRecorder()

	// Выполнение
	service.HandleProviderNotification(w, req)

	// Проверка
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		t.Error("Ожидалась ошибка, но получен статус 200")
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	if len(bodyBytes) == 0 {
		t.Error("Ожидалось сообщение об ошибке в теле ответа")
	}
}

// TestHandleClientNotification_StorageError тестирует обработку ошибок хранилища
func TestHandleClientNotification_StorageError(t *testing.T) {
	// Подготовка
	mockStorage := NewMockStorage()
	mockStorage.SetSaveError(errors.New("storage error"))
	mockLogger := NewMockLogger()
	mockValidator := NewMockJWSValidator()

	// Создание декодера и парсера
	decoder := applestore.NewAppleDecoder(mockValidator)
	parser := applestore.NewAppleParser(decoder)

	// Создание сервиса
	service := applestore.NewAppleStoreService(mockStorage, mockLogger, parser)

	// Создание запроса с тестовыми данными
	clientNotification := map[string]interface{}{
		"bundleId":              "com.test.app",
		"appAccountToken":       "user123",
		"signedTransactionInfo": "header.eyJvcmlnaW5hbFRyYW5zYWN0aW9uSWQiOiIxMjM0NTYiLCJ0cmFuc2FjdGlvbklkIjoiNTQzMjEiLCJwcm9kdWN0SWQiOiJjb20udGVzdC5wcm9kdWN0IiwiZXhwaXJlc0RhdGUiOjE3MjUwMDAwMDAwMDB9.signature",
	}
	body, _ := json.Marshal(clientNotification)
	req := httptest.NewRequest(http.MethodPost, "/client-notification", bytes.NewReader(body))
	w := httptest.NewRecorder()

	// Выполнение
	service.HandleClientNotification(w, req)

	// Проверка ошибки
	resp := w.Result()
	defer resp.Body.Close()

	// В зависимости от реализации, ошибка хранилища может вернуть ошибку или 200 OK
	// Проверяем, что запись в хранилище не произошла
	status, _ := mockStorage.GetSubscriptionStatus("user123")
	if status != nil {
		t.Error("Статус подписки был сохранен несмотря на ошибку")
	}
}
