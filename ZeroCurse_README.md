# ZeroCurse AI Antivirus 🛡️🌀

A proactive WordPress plugin that detects and limits unauthorized AI scrapers (e.g., GPTBot, Claude-bot, Google-Extended) using a Tar Pit (data trap) mechanism.

## How it works

Unlike traditional firewalls that simply block IPs, ZeroCurse makes scraping more resource‑intensive for the bot by increasing the volume of data it must process.

1. **Detection.** Identifies AI crawlers via User‑Agent and other signatures at the earliest server load stage (using the `plugins_loaded` hook).
2. **Isolation.** Immediately blocks access to real site content by returning HTTP 403.
3. **Active Trap (Tar Pit).** Sends a stream of JSON chunks to the bot. This increases the amount of data the scraper must download and parse, raising its computational and bandwidth costs.

> **Important:** Effectiveness depends on how the specific scraper handles large responses and unusual numeric patterns. Some clients may ignore the payload or terminate the connection early upon detecting anomalies.

## Technical details

- **Language:** PHP (compatible with standard WordPress core).
- **Defense method:** Active Defense / Tar Pit.
- **Trigger:** `plugins_loaded` hook for early request interception.
- **Server safety:** Built‑in loop limiter and output volume control protect the host server from overload.

## Known limitations

- The mathematical pattern (`0/0` raised to `-0`) does not guarantee crashes across all parsers; behavior depends on the scraper’s tokenizer and error‑handling logic.
- Advanced crawlers may detect abnormally large responses and abort the request early.
- This plugin increases outbound traffic for blocked bots; consider hosting bandwidth limits.
- Not intended to block legitimate search engine bots; configure exclusions before enabling active defense.

## Quick start

1. Download the plugin ZIP from the latest release.
2. In WordPress admin: Plugins → Add New → Upload Plugin.
3. Activate **ZeroCurse AI Antivirus**.
4. Go to Settings → ZeroCurse and toggle **Active Defense** to ON.
5. (Optional) Add custom User‑Agent patterns to the blocklist or exclusions.

## License

MIT License. Free for personal and commercial use.


# ZeroCurse AI Antivirus 🛡️🌀

Проактивный плагин для WordPress, который выявляет и ограничивает работу неавторизованных ИИ‑скраперов (например, GPTBot, Claude‑bot, Google‑Extended) с помощью механизма Tar Pit («дегтярная яма»).

## Как это работает

В отличие от традиционных файрволов, которые просто блокируют IP‑адреса, ZeroCurse делает скрапинг дороже для бота за счёт увеличения объёма передаваемых данных и сложности их обработки.

1. **Обнаружение.** Плагин распознаёт ИИ‑краулеры по User‑Agent и другим сигнатурам на ранней стадии загрузки сервера (через хук `plugins_loaded`).
2. **Изоляция.** Доступ к реальному контенту сайта немедленно блокируется: боту возвращается HTTP‑код 403.
3. **Активная ловушка (Tar Pit).** Вместо пустого ответа плагин отдаёт поток JSON‑чанков. Это увеличивает объём данных, которые бот должен скачать и обработать, и повышает затраты его вычислительных ресурсов.

> **Важно:** Эффективность зависит от того, как конкретный бот обрабатывает нестандартные данные и большие ответы. Некоторые парсеры могут игнорировать JSON или сразу прерывать соединение при обнаружении аномалий.

## Технические детали

- **Язык:** PHP (совместим со стандартным ядром WordPress).
- **Метод защиты:** активная оборона / Tar Pit.
- **Триггер:** хук `plugins_loaded` для максимально раннего перехвата запроса.
- **Безопасность сервера:** встроенный ограничитель циклов и контроль объёма отдачи данных защищают сервер сайта от перегрузки.

## Известные ограничения

- Математическая конструкция (`0/0` в степени `-0`) не гарантирует сбой всех парсеров: поведение зависит от реализации токенизатора
