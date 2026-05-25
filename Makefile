# Makefile для проекта push-swap

.PHONY: all clean checker push-swap test

# Сборка обоих программ (checker первым по требованиям)
all: checker push-swap

# checker — проверяет корректность инструкций
checker:
	go build -o checker ./cmd/checker

# push-swap — генерирует инструкции сортировки
push-swap:
	go build -o push-swap ./cmd/push-swap

# Очистка
clean:
	rm -f checker checker.exe push-swap push-swap.exe

# Запуск тестов
test:
	go test ./...
