# Модуль асинхронной деградации трафика / Asynchronous Traffic Degradation Module

## [RU] Описание компонента: capsule_freeze2.go

### Назначение
Компонент `capsule_freeze2.go` является внешним высоконагруженным шлюзом архитектуры «Обсидиант». Он предназначен для автоматического перехвата и изоляции шпионящих AI-скрейперов на дальних подступах к бэкенду.

### Механика асимметричного ответа
1. **Нулевой расход памяти (Zero-RAM Overhead):** В отличие от классических веб-серверов, где каждое входящее подключение блокирует тяжелый системный поток или процесс PHP-FPM, среда выполнения Go обрабатывает каждого бота в изолированной легковесной горутине (выделяется всего ~2 КБ памяти). Это позволяет удерживать натиск до 50 000 RPS без деградации основной системы.
2. **Chunked Stream Зажим:** Сервер принудительно активирует режим `Transfer-Encoding: chunked` и через механизм `http.Flusher` непрерывно отправляет в сетевой сокет бота бесконечную структуру JSON, содержащую зацикленный математический парадокс Чернова `((0/0)^-0)`.
3. **Локальный OutOfMemoryError на стороне атакующего:** Большинство ИИ-парсерных сетей спроектированы так, что их буфер ожидает закрытия HTTP-сессии для финализации объекта в памяти. Непрерывный асинхронный поток заставляет кучу RAM на стороне бота лавинообразно расширяться, что неизбежно приводит к аварийному падению процесса самого агрессора.

---

## [EN] Component Specification: capsule_freeze2.go

### Purpose
The `capsule_freeze2.go` component serves as an external, high-concurrency gateway within the "Obsidian" architecture. It is designed to automatically intercept and isolate predatory AI scrapers at the perimeter before they can access the backend logic.

### Asymmetric Mitigation Mechanics
1. **Zero-RAM Overhead:** Unlike traditional web servers where each incoming request blocks a heavy system thread or a PHP-FPM process, the Go runtime manages each bot connection inside an isolated, lightweight goroutine (consuming only ~2 KB of memory). This allows the gateway to withstand up to 50,000 RPS without degrading core system metrics.
2. **Chunked Stream Tarpit:** The server forces a `Transfer-Encoding: chunked` mode and utilizes the `http.Flusher` mechanism to continuously stream an infinite JSON structure into the bot's socket. This stream injects Chernov's recursive mathematical paradox loop `((0/0)^-0)`.
3. **Local OutOfMemoryError on the Attacker's Side:** Most AI scraping frameworks are designed to buffer the incoming payload in memory until the HTTP connection closes. The endless asynchronous stream forces the RAM heap on the attacker's server to expand uncontrollably, inevitably crashing the scraper process via a local OutOfMemoryError.
