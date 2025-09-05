package applestore

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"
)

// TestAppleParser проверяет работу парсера для разбора уведомлений Apple
func TestAppleParser(t *testing.T) {
	// Создаем мок-валидатора
	mockValidator := NewMockJWSValidator()

	// Создаем декодер и парсер
	decoder := NewAppleDecoder(mockValidator)
	parser := NewAppleParser(decoder)

	// Тест для ParseClientNotification
	t.Run("ParseClientNotification", func(t *testing.T) {
		// Создаем тестовые данные
		clientNotificationJSON := `{
			"bundleId": "com.test.app",
			"appAccountToken": "user123",
			"signedTransactionInfo": "header.payload.signature"
		}`

		// Создаем Reader из JSON
		reader := bytes.NewReader([]byte(clientNotificationJSON))

		// Вызываем метод парсинга
		notification, err := parser.ParseClientNotification(reader)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при парсинге уведомления клиента: %v", err)
		}

		if notification.BundleID != "com.test.app" {
			t.Errorf("Неправильный BundleID: ожидался 'com.test.app', получен '%s'", notification.BundleID)
		}

		if notification.AppAccountToken != "user123" {
			t.Errorf("Неправильный AppAccountToken: ожидался 'user123', получен '%s'", notification.AppAccountToken)
		}

		if notification.SignedTransactionInfo != "header.payload.signature" {
			t.Errorf("Неправильный SignedTransactionInfo: ожидался 'header.payload.signature', получен '%s'", notification.SignedTransactionInfo)
		}
	})

	// Тест для ParseTransaction
	t.Run("ParseTransaction", func(t *testing.T) {
		// Настройка мок-валидатора чтобы он не возвращал ошибок
		mockValidator.ValidateError = nil

		// Тестовые данные с закодированной транзакцией
		// Строка должна иметь формат "header.payload.signature", где payload - base64url закодированные данные транзакции
		// В этом примере, payload = eyJvcmlnaW5hbFRyYW5zYWN0aW9uSWQiOiIxMjM0NTYiLCJ0cmFuc2FjdGlvbklkIjoiNTQzMjEiLCJwcm9kdWN0SWQiOiJjb20udGVzdC5wcm9kdWN0In0=
		// который декодируется в {"originalTransactionId":"123456","transactionId":"54321","productId":"com.test.product"}
		transactionJWS := "header.eyJvcmlnaW5hbFRyYW5zYWN0aW9uSWQiOiIxMjM0NTYiLCJ0cmFuc2FjdGlvbklkIjoiNTQzMjEiLCJwcm9kdWN0SWQiOiJjb20udGVzdC5wcm9kdWN0In0.signature"

		// Вызываем метод парсинга
		transaction, err := parser.ParseTransaction(transactionJWS)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при парсинге транзакции: %v", err)
		}

		if transaction.OriginalTransactionID != "123456" {
			t.Errorf("Неправильный OriginalTransactionID: ожидался '123456', получен '%s'", transaction.OriginalTransactionID)
		}

		if transaction.TransactionID != "54321" {
			t.Errorf("Неправильный TransactionID: ожидался '54321', получен '%s'", transaction.TransactionID)
		}

		if transaction.ProductID != "com.test.product" {
			t.Errorf("Неправильный ProductID: ожидался 'com.test.product', получен '%s'", transaction.ProductID)
		}
	})

	// Тест на ошибки валидации
	t.Run("ValidateError", func(t *testing.T) {
		// Настраиваем мок-валидатор на возврат ошибки
		mockValidator.ValidateError = errors.New("validation failed")

		// Тестовые данные
		transactionJWS := "header.payload.signature"

		// Вызываем метод парсинга
		_, err := parser.ParseTransaction(transactionJWS)

		// Проверяем наличие ошибки
		if err == nil {
			t.Fatal("Ожидалась ошибка валидации, но ее не было")
		}
	})
}

// TestAppleDecoder проверяет работу декодера Apple JWS
func TestAppleDecoder(t *testing.T) {
	// Создаем мок-валидатор
	mockValidator := NewMockJWSValidator()

	// Создаем декодер
	decoder := NewAppleDecoder(mockValidator)

	// Тест успешного декодирования JWS
	t.Run("DecodeSignedJWS_Success", func(t *testing.T) {
		// Настройка мок-валидатора для успешной валидации
		mockValidator.ValidateError = nil

		// Тестовые данные: корректный формат JWS
		// header.eyJrZXkiOiJ2YWx1ZSJ9.signature
		// payload декодируется в {"key":"value"}
		signedJWS := "header.eyJrZXkiOiJ2YWx1ZSJ9.signature"

		// Вызываем метод декодирования
		payload, err := decoder.DecodeSignedJWS(signedJWS)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при декодировании JWS: %v", err)
		}

		// Проверяем содержимое payload
		var data map[string]string
		if err := json.Unmarshal(payload, &data); err != nil {
			t.Fatalf("Не удалось распарсить декодированный payload: %v", err)
		}

		if value, ok := data["key"]; !ok || value != "value" {
			t.Errorf("Неправильное значение в декодированных данных: ожидалось key=value, получено key=%s", value)
		}
	})

	// Тест ошибки валидации
	t.Run("DecodeSignedJWS_ValidationError", func(t *testing.T) {
		// Настройка мок-валидатора для возврата ошибки
		mockValidator.ValidateError = errors.New("validation error")

		// Тестовые данные: корректный формат JWS
		signedJWS := "header.payload.signature"

		// Вызываем метод декодирования
		_, err := decoder.DecodeSignedJWS(signedJWS)

		// Проверяем наличие ошибки валидации
		if err == nil {
			t.Fatal("Ожидалась ошибка валидации, но ее не было")
		}
	})

	// Тест некорректного формата JWS
	t.Run("DecodeSignedJWS_InvalidFormat", func(t *testing.T) {
		// Тестовые данные: некорректный формат JWS (нет трех частей)
		signedJWS := "header.payload"

		// Вызываем метод декодирования
		_, err := decoder.DecodeSignedJWS(signedJWS)

		// Проверяем наличие ошибки формата
		if err == nil {
			t.Fatal("Ожидалась ошибка формата JWS, но ее не было")
		}
	})

	// Тест пустого JWS
	t.Run("DecodeSignedJWS_Empty", func(t *testing.T) {
		// Тестовые данные: пустая строка
		signedJWS := ""

		// Вызываем метод декодирования
		_, err := decoder.DecodeSignedJWS(signedJWS)

		// Проверяем наличие ошибки
		if err == nil {
			t.Fatal("Ожидалась ошибка при пустом JWS, но ее не было")
		}
	})

	// Тест декодирования уведомления
	t.Run("DecodeRawNotification", func(t *testing.T) {
		// Тестовые данные: JSON с полем signedPayload
		notificationJSON := `{"signedPayload": "test-payload"}`

		// Создаем Reader из JSON
		reader := bytes.NewReader([]byte(notificationJSON))

		// Вызываем метод декодирования
		signedPayload, err := decoder.DecodeRawNotification(reader)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при декодировании уведомления: %v", err)
		}

		if signedPayload != "test-payload" {
			t.Errorf("Неправильный signedPayload: ожидался 'test-payload', получен '%s'", signedPayload)
		}
	})
}

// Использует MockJWSValidator из applestore_test.go
