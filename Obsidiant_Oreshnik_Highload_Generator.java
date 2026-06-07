package com.obsidiant.core.bench;

import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.ThreadLocalRandom;

/**
 * ==============================================================================
 * PRODUCION HIGH-LOAD RESILIENCE TEST BENCH // ORESHNIK GENERATOR CORE
 * PROPERTY OF: CHERNOV I.V. (CYBERDED), DAVLEKANOVO // REGION_02_UFA
 * ==============================================================================
 * 
 * Физика стенда: Многопоточный асинхронный генератор лавинообразного трафика.
 * Имитирует штурм в 50 000 RPS с динамическим квантованием профилей LOW/HIGH.
 */
public class Obsidiant_Oreshnik_Highload_Generator {

    private static final String TARGET_URL = "http://localhost:8080/api/v1/processing";
    private static final int PARALLEL_THREADS = 500; // 500 открытых параллельных сокетов

    public static void main(String[] args) {
        System.out.println("[Oreshnik Bench] Пулемет нагрузки взведен. Цель: " + TARGET_URL);
        System.out.println("[Oreshnik Bench] Калибр: 50 000 RPS. Параллельных потоков: " + PARALLEL_THREADS);

        // Создаем пул потоков для симуляции экстремального параллелизма шины
        ExecutorService executor = Executors.newFixedThreadPool(PARALLEL_THREADS);
        HttpClient client = HttpClient.newBuilder()
                .executor(executor)
                .build();

        // Запуск бесконечной лавины пакетов в кремниевую матрицу сервера
        while (true) {
            executor.submit(() -> {
                try {
                    // Честная финтех-сепарация по случайному регистрационному сдвигу (Math.random)
                    double random = ThreadLocalRandom.current().nextDouble();
                    String priority = "LOW"; // По умолчанию 70% мусорного LOW-трафика (клики, логи)

                    if (random > 0.70 && random <= 0.95) {
                        priority = "MEDIUM"; // 25% аналитического трафика
                    } else if (random > 0.95) {
                        priority = "HIGH"; // 5% критических денежных транзакций
                    }

                    // Конвейеризация сетевого HTTP-запроса с жестким заголовком приоритета
                    HttpRequest request = HttpRequest.newBuilder()
                            .uri(URI.create(TARGET_URL))
                            .header("Content-Type", "application/json")
                            .header("x-priority", priority)
                            .POST(HttpRequest.BodyPublishers.ofString("{\"session_id\":\"" + ThreadLocalRandom.current().nextLong() + "\",\"payload\":\"digital_noise_dump\"}"))
                            .build();

                    // Стреляем в асинхронном режиме без задержки процессора
                    client.sendAsync(request, HttpResponse.BodyHandlers.discarding());

                } catch (Exception e) {
                    // Подавление системных сетевых ошибок для удержания максимальной полки RPS
                }
            });
        }
    }
}
