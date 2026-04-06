#!/bin/sh

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Значения по умолчанию
DEFAULT_TARGET_DIR="tests/cases"
DEFAULT_EXECUTABLE_PATH="bin/fractalflame"

# Директория, переданная в качестве аргумента, или значение по умолчанию
if [ -z "$1" ]; then
    TARGET_DIR="$DEFAULT_TARGET_DIR"
    echo "${YELLOW}Используется директория по умолчанию: $TARGET_DIR${NC}"
else
    TARGET_DIR="$1"
fi

# Путь к исполняемому файлу, или значение по умолчанию
if [ -z "$2" ]; then
    EXECUTABLE_PATH="$DEFAULT_EXECUTABLE_PATH"
    echo "${YELLOW}Используется исполняемый файл по умолчанию: $EXECUTABLE_PATH${NC}"
    # Проверяем, существует ли бинарник, если нет - пытаемся собрать
    if [ ! -f "$EXECUTABLE_PATH" ]; then
        echo "${YELLOW}Бинарник не найден, пытаемся собрать...${NC}"
        go build -o "$EXECUTABLE_PATH" ./cmd/fractalflame
        if [ $? -ne 0 ]; then
            echo "${RED}Ошибка: Не удалось собрать бинарник.${NC}"
            echo "${YELLOW}Использование: $0 [<путь_к_директории>] [<путь_к_исполняемому_файлу>]${NC}"
            exit 1
        fi
    fi
else
    EXECUTABLE_PATH="$2"
fi

# Проверяем, существует ли указанная директория
if [ ! -d "$TARGET_DIR" ]; then
    echo "${RED}Ошибка: Директория $TARGET_DIR не существует.${NC}"
    exit 1
fi

# Проверяем, существует ли исполняемый файл
if [ ! -f "$EXECUTABLE_PATH" ]; then
    echo "${RED}Ошибка: Исполняемый файл $EXECUTABLE_PATH не найден.${NC}"
    exit 1
fi

# Делаем файл исполняемым (на случай, если он не имеет прав на выполнение)
chmod +x "$EXECUTABLE_PATH"

# Переменные для подсчета тестов
total_tests=0
successful_tests=0

# Перебираем все .sh файлы в указанной директории
echo "${BLUE}Ищем все .sh файлы в директории: $TARGET_DIR${NC}"

for script in "$TARGET_DIR"/*.sh; do
    # Проверяем, существует ли файл
    if [ -f "$script" ]; then
        total_tests=$((total_tests + 1))
        echo "${YELLOW}Запускаем: $script${NC}"
        # Делаем файл исполняемым (на случай, если он не имеет прав на выполнение)
        chmod +x "$script"
        # Запускаем файл
        sh "$script" "$EXECUTABLE_PATH"
        # Проверяем код возврата
        if [ $? -ne 0 ]; then
            echo "${RED}Ошибка при выполнении $script${NC}"
        else
            echo "${GREEN}$script выполнен успешно${NC}"
            successful_tests=$((successful_tests + 1))
        fi
    else
        echo "${RED}Файл $script не найден или не является исполняемым${NC}"
    fi
done

# Вывод результатов
echo ""
echo "${BLUE}Результаты тестирования:${NC}"
echo "Всего тестов: $total_tests"
echo "Успешных тестов: $successful_tests"

# Проверяем, было ли выполнено ровно 3 теста
if [ "$total_tests" -ne 3 ]; then
    echo "${RED}Ошибка: Ожидалось 3 теста, но выполнено $total_tests.${NC}"
    exit 1
fi

# Если все тесты успешны
if [ "$total_tests" -eq "$successful_tests" ]; then
    echo "${GREEN}Все тесты успешно пройдены!${NC}"
    exit 0
else
    echo "${RED}Некоторые тесты завершились с ошибками.${NC}"
    exit 1
fi
