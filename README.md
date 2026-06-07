# 📐 obsidiant-64gb-oreshnik-bench

### 🇷🇺 Испытательный полигон «Обсидиант: Орешник» для банковских систем 64 ГБ+
**Добро пожаловать на полигон КиберДеда!** Этот проект создан для жестких краш-тестов, симуляции лавинных перегрузок и ликвидации аварий в банковских микросервисах с объемом памяти **64 ГБ и выше** без снятия тяжелых бинарных дампов (`HPROF`).

**Компоненты полигона:**
1. **💥 Пулеметный стенд «Орешник» (`obsidiant-test-bench.js`):** Имитатор экстремального хайлоада на Node.js. Работает в режиме непрерывной стрельбы, засыпая память лавиной транзакций и генерируя до **75% обломков и мусора**, искусственно доводя сервер до критической точки падения.
2. **🚀 Модуль «Обратный Трассировщик» (`Obsidiant_64GB_Oreshnik_Proving_Ground_Agent.java`):** Интеллектуальный Java-агент. Сидит в засаде внутри JVM. Как только «Орешник» забивает память до **85%**, агент блокирует OOMKill и за 5 секунд раскручивает граф ссылок назад от зомби-объектов к корням (`GC Roots`), выгружая в Grafana легкую текстовую схему в пару мегабайт вместо 64 ГБ тяжелого дампа. Безопасность ИБ соблюдена на 100%.

---

### 🇺🇸 "Obsidiant: Oreshnik" Proving Ground for 64 GB+ Banking Systems
**Welcome to CyberDed's Proving Ground!** This repository is designed for heavy crash-testing, simulating avalanche-like traffic overloads, and eliminating memory failures in enterprise banking microservices with **64 GB+ RAM** without triggering destructive binary heap dumps (`HPROF`).

**Proving Ground Architecture:**
1. **💥 Machine-Gun "Oreshnik" (`obsidiant-test-bench.js`):** An extreme high-throughput load emulator on Node.js. It operates in a rapid-fire machine-gun mode, flooding the system with synthetic transactions, generating up to **75% of residual memory debris and garbage**, driving the server to its breaking point.
2. **🚀 "Reverse Tracer" Module (`Obsidiant_64GB_Oreshnik_Proving_Ground_Agent.java`):** A smart Java Agent embedded within the JVM. Once the "Oreshnik" generator fills the memory up to **85%**, the agent blocks the OOMKill and traces object references backward to `GC Roots` in under 5 seconds, sending a lightweight text tree (a few megabytes) to Grafana instead of a massive 64 GB dump. 100% InfoSec compliant.

---
📐 **io.github.ssmel63.obsidiant** | Engineered by CyberDed / КиберДед (Игорь Чернов)
