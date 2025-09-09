package applestore

import (
	"context"
	"errors"
	"testing"

	"subscription-server/internal/applestore"
	"subscription-server/internal/contracts"
	"subscription-server/internal/logger"
	"subscription-server/internal/storage"
)

// TestMocks проверяет корректность работы моков
func TestMocks(t *testing.T) {
	// Тест мока хранилища
	t.Run("MockStorage", func(t *testing.T) {
		mock := NewMockStorage()

		// Проверяем, что начальное состояние пусто
		status, err := mock.GetSubscriptionStatus(context.Background(), "testUser")
		if err != nil {
			t.Errorf("Получена ошибка при первом запросе GetSubscriptionStatus: %v", err)
		}
		if status != nil {
			t.Errorf("Начальное состояние должно быть nil, получено %+v", status)
		}

		// Проверяем установку ошибки
		testError := errors.New("test error")
		mock.SetGetError(testError)
		_, err = mock.GetSubscriptionStatus(context.Background(), "testUser")
		if err != testError {
			t.Errorf("Ожидалась ошибка %v, получена %v", testError, err)
		}

		// Сбрасываем ошибку и проверяем сохранение
		mock.SetGetError(nil)
		testStatus := &storage.SubscriptionStatus{UserToken: "testUser", ProductID: "testProduct"}
		err = mock.SetSubscriptionStatus(context.Background(), testStatus)
		if err != nil {
			t.Errorf("Ошибка при сохранении статуса: %v", err)
		}

		// Проверяем, что данные сохранились
		savedStatus, err := mock.GetSubscriptionStatus(context.Background(), "testUser")
		if err != nil {
			t.Errorf("Ошибка при получении сохраненного статуса: %v", err)
		}
		if savedStatus == nil {
			t.Fatal("Сохраненный статус не найден")
		}
		if savedStatus.UserToken != "testUser" || savedStatus.ProductID != "testProduct" {
			t.Errorf("Некорректный сохраненный статус: %+v", savedStatus)
		}
	})

	// Тест мока логгера
	t.Run("MockLogger", func(t *testing.T) {
		mock := NewMockLogger()

		// Проверяем, что начальное состояние пусто
		if len(mock.logs) > 0 {
			t.Errorf("Начальное состояние логов должно быть пустым, получено %d записей", len(mock.logs))
		}

		// Добавляем лог и проверяем
		testMessage := logger.LogMessage{
			Level:   "INFO",
			Sender:  "Test",
			Message: "Test message",
		}
		mock.Log(testMessage)

		if len(mock.logs) != 1 {
			t.Fatalf("Ожидалась 1 запись в логах, получено %d", len(mock.logs))
		}

		loggedMessage := mock.logs[0]
		if loggedMessage.Level != "INFO" || loggedMessage.Sender != "Test" || loggedMessage.Message != "Test message" {
			t.Errorf("Некорректное сообщение в логе: %+v", loggedMessage)
		}
	})

	// Тест мока валидатора
	t.Run("MockJWSValidator", func(t *testing.T) {
		mock := NewMockJWSValidator()

		// По умолчанию ошибки нет
		err := mock.Validate("header", "payload", "signature")
		if err != nil {
			t.Errorf("Неожиданная ошибка при валидации: %v", err)
		}

		// Устанавливаем ошибку и проверяем
		testError := errors.New("validation failed")
		mock.SetValidateError(testError)

		err = mock.Validate("header", "payload", "signature")
		if err != testError {
			t.Errorf("Ожидалась ошибка %v, получена %v", testError, err)
		}
	})
}

// TestFactory проверяет правильность создания сервиса и компонентов через фабрики
func TestFactory(t *testing.T) {
	// Тест создания декодера
	t.Run("NewAppleDecoder", func(t *testing.T) {
		validator := NewMockJWSValidator()
		decoder := applestore.NewAppleDecoder(validator)

		if decoder == nil {
			t.Fatal("Декодер не был создан")
		}
	})

	// Тест создания парсера
	t.Run("NewAppleParser", func(t *testing.T) {
		validator := NewMockJWSValidator()
		decoder := applestore.NewAppleDecoder(validator)
		parser := applestore.NewAppleParser(decoder)

		if parser == nil {
			t.Fatal("Парсер не был создан")
		}
	})

	// Тест создания сервиса
	t.Run("NewAppleStoreService", func(t *testing.T) {
		storage := NewMockStorage()
		logger := NewMockLogger()
		validator := NewMockJWSValidator()
		decoder := applestore.NewAppleDecoder(validator)
		parser := applestore.NewAppleParser(decoder)

		service := applestore.NewAppleStoreService(storage, logger, parser)

		// Проверяем, что сервис создан и реализует интерфейс contracts.Service
		var _ contracts.Service = service

		if service == nil {
			t.Fatal("Сервис не был создан")
		}
	})
}
