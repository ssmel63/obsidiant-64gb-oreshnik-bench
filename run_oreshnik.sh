#!/bin/bash
# ==============================================================================
# JAVA LAUNCHER SCRIPT // PROJECT: OBSIDIANT ORESHNIK RESILIENCE BENCH
# PROPERTY OF: CHERNOV I.V. (CYBERDED), DAVLEKANOVO // REGION_02_UFA
# ==============================================================================

echo "[Oreshnik Java] Запуск автоматической компиляции генератора шторма..."

# Проверяем физическое наличие компилятора Java в операционной системе Linux банка
if ! command -v javac &> /dev/null; then
    echo "ОШИБКА: Компилятор Java (javac) не найден в системе Linux банка."
    exit 1
fi

# Шаг 1. Компилируем Java-файл нашего пулемета нагрузки
javac -d . Obsidiant_Oreshnik_Highload_Generator.java

if [ $? -eq 0 ]; then
    echo "[Oreshnik Java] Компиляция завершена успешно. Байт-код JVM собран."
    echo "[Oreshnik Java] ВНИМАНИЕ: Активация лавины на 50 000 RPS. Для остановки нажмите Ctrl + C"
    
    # Шаг 2. Запускаем скомпилированный класс внутри виртуальной машины Java
    java com.obsidiant.core.bench.Obsidiant_Oreshnik_Highload_Generator
else
    echo "ОШИБКА: Сбой компиляции Java-класса. Проверьте исходный код."
    exit 1
fi
