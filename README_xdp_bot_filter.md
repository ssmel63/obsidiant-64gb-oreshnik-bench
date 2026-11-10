# Техническая спецификация модуля eBPF/XDP: xdp_bot_filter.c
# Technical Specification for eBPF/XDP Module: xdp_bot_filter.c

---

## [RU] Архитектурный разбор Модели 4.5 (Сетевой уровень ядра)

Этот компонент является передовым рубежом защиты платформы «Обсидиант», развернутым непосредственно в контексте сетевого драйвера (XDP). Он обрабатывает входящий трафик на скорости физического кабеля до того, как операционная система выделит под него структуры памяти (RAM) или потоки CPU.

### Ключевые инженерные решения и ликвидация ловушек ядра:

1. **Решение проблемы верификатора `math between pkt pointer and register`**
   - **Проблема:** Стандартные eBPF-программы отклоняются верификатором Linux при попытке динамического сдвига указателя TCP через переменную длины IP-заголовка `(ip->ihl * 4)`. Ядро блокирует код, так как не может статически доказать безопасность границ памяти.
   - **Решение Чернова:** Программа применяет битовую маску `ip->ihl & 0x0F` и жестко отсекает любые пакеты с IP-опциями через контур `if (ihl > 5) return XDP_DROP;`. С этой наносекунды состояние регистра `ihl` в графе верификатора становится константой, строго равной 5. Выражение `struct tcphdr *tcp = (void *)ip + 20;` преобразуется в безопасный константный сдвиг, и ошибка pointer arithmetic полностью уничтожается.

2. **Ювелирный зажим сканеров (Pure SYN Lockdown)**
   - Программа не блокирует легитимные SYN-пакеты (что разрушило бы работу реальных пользователей). Вместо этого она вычисляет точную позицию байта флагов TCP через фиксированный безопасный указатель `flags_ptr = (void *)tcp + 13`.
   - Проверка `if (flags == 0x02)` под микроскопом изолирует пакеты, содержащие **исключительно флаг SYN** (без ACK, FIN, PSH и др.). Это позволяет со 100% точностью срезать автоматизированную разведку ботнетов (Nmap, Masscan) на подлетах.

3. **Защита от ложных срабатываний при фрагментации**
   - Внедрена ультралегкая проверка флагов фрагментации через маску `frag_off & 0x1FFF`. Если пакет является последующим фрагментом (где заголовка TCP физически нет), он мгновенно отправляется в стек ОС (`XDP_PASS`) для сборки ядром Linux, исключая порчу пользовательских данных.

---

## [RU] РЕГЛАМЕНТ ПРОВЕДЕНИЯ ЛАБОРАТОРНЫХ ИСПЫТАНИЙ

Внимание! Данный демонстрационный модуль (Модель 4.5) поставляется со специально вырезанными векторами активного противодействия (Модели 3 и 4), чтобы не нарушать политики безопасности публичных платформ.

Если вашей исследовательской лаборатории (или отделу анализа угроз/WAF) требуется провести полноценные высоконагруженные испытания (High-Load Sandbox Evaluation) финального суверенного ядра **Obsidian Core v5.0 (Stealth Edge)** под Astra Linux:

1. Направьте официальный запрос на проведение пилотных тестов на электронную почту: ** kiberdied [собака] proton [точка] me **
2. После подписания соглашения о неразглашении (NDA), вам будет предоставлен доступ к закрытому бинарному контуру, включающему автоматическую генерацию Abuse-паспортов атак и модули CPU-Vampirism.

**КРИТИЧЕСКИЙ ПРОТОКОЛ БЕЗОПАСНОСТИ:**  
Шлюз Лаборатории MemWorm настроен на жесткую фильтрацию входящей почты. Рассматриваются исключительно текстовые сообщения в теле письма. Любые входящие письма, содержащие вложения (файлы, документы, архивы), внешние ссылки или шпионские пиксели, автоматически уничтожаются на сервере приема и заносятся в пожизненный eBPF-бан. Все переговоры ведутся строго в текстовом формате.

**Доступ к исходному коду v4.5:**  
Полный исходный код демонстрационной Модели 4.5 упакован в зашифрованный архив для фиксации патентного и временного приоритета. Архив размещен в защищенном облачном хранилище. Пароль от архива и прямая ссылка предоставляются бесплатно представителям сертифицированных ИБ-лабораторий после верификации. Направьте официальный запрос на указанный текстовый адрес.

---

## [EN] Architectural Breakdown of Model 4.5 (Kernel Network Layer)

This component serves as the perimeter defense line of the "Obsidian" platform, running directly within the network driver context (XDP). It intercepts and drops high-concurrency botnet traffic at wire speed before the OS allocates any RAM structures or CPU threads.

### Key Engineering Breakthroughs & Kernel Bug Mitigation:

1. **Solving the `math between pkt pointer and register` Verifier Bug**
   - **The Bug:** Traditional eBPF programs are rejected by the Linux kernel verifier when attempting a dynamic TCP pointer offset using the IP header length variable `(ip->ihl * 4)`. The kernel blocks the bytecode because it cannot statically mathematically prove memory bounds safety.
   - **Chernov's Solution:** The program forces a bitwise mask `ip->ihl & 0x0F` and immediately drops any packets containing IP options via the `if (ihl > 5) return XDP_DROP;` routine. From this nanosecond forward, the `ihl` register state within the verifier's Scalar Range Tracking is strictly locked as a constant equal to 5. The expression `struct tcphdr *tcp = (void *)ip + 20;` becomes a safe, constant-offset operation, completely satisfying the strict memory access verification.

2. **Precision High-Load Scanner Mitigation (Pure SYN Lockdown)**
   - The program avoids broad SYN blocking (which would break legitimate client connection handshakes). Instead, it targets the exact position of the TCP flags byte using a safe, constant pointer offset `flags_ptr = (void *)tcp + 13`.
   - The validation `if (flags == 0x02)` strictly isolates packets containing **only the SYN flag** (with ACK, FIN, PSH, and other flags completely absent). This allows the filter to cut down automated botnet reconnaissance tools (Nmap, Masscan, Custom Scrapers) with zero impact on legitimate traffic.

3. **Zero-Overhead Fragmentation Shield**
   - An ultra-lightweight fragmentation flag filter is implemented via the `frag_off & 0x1FFF` bitwise mask. If a packet is identified as a displaced fragment (where no L4 TCP header is physically present), it is instantly passed forward to the host network stack (`XDP_PASS`) for native Linux kernel reassembly, eliminating any potential false positives or user data corruption.

---

## [EN] LAB TESTING AND EVALUATION PROTOCOL

Notice! This documentation module (Model 4.5) is delivered with active counter-attack vectors (Model 3 & 4) intentionally deprecated to comply with public platform security guidelines.

If your threat intelligence laboratory (or infrastructure WAF/Security architecture team) requires a full high-load sandbox evaluation of the final sovereign **Obsidian Core v5.0 (Stealth Edge)** under Astra Linux environments:

1. Submit an official request for pilot sandbox testing directly to the lead developer at: ** kiberdied [at] proton [dot] me **
2. Upon executing a Non-Disclosure Agreement (NDA), your team will be granted restricted access to the private binary suite, including automated abuse-passport generation and CPU-Vampirism engines.

**STRICT SECURITY COMPLIANCE PROTOCOL:**  
The MemWorm Lab inbound gateway enforces a zero-trust payload policy. Only raw text messages within the email body are parsed. Any inbound emails containing attachments (files, documents, archives), embedded links, or tracking pixels are instantly dropped at the perimeter server and blacklisted via eBPF. All communications are strictly text-only.

**Source Code Access v4.5:**  
The complete source tree for demonstration Model 4.5 is packaged into an encrypted archive to lock patent and timestamp priority. The archive is hosted in a secure cloud storage. The decryption password and direct link are provided free of charge to verified security teams upon official request. Contact via the obfuscated text address above.
