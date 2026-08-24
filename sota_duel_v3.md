# 🌐 Sota-Junior / Чернов-Core v3: «Принцип Брюквы» (The Rutabaga Principle)

[RU] Международная патентно-архитектурная спецификация асинхронного финтех-ядра для детской интерактивной геймификации.
[EN] International patent-architectural specification of the asynchronous fintech core for interactive children's gamification.

---

## [RU] РУССКАЯ ВЕРСИЯ

### 🎯 Суть концепции («Принцип Брюквы»)
Традиционные банковские инвест-копилки для взрослых скучны и неповоротливы. Данное решение переводит автоматическое округление повседневных трат родителя (например, покупка брюквы или кваса с 555 до 560 рублей) в плоскость эмоционального детского гейминга. Продукт формирует у подрастающего поколения лояльность к банку на уровне «грудного детского молока», привязывая миллионы будущих клиентов через теплые воспоминания.

### 🛠 Проблема ИТ-инфраструктуры банков
Если миллионы детей начнут азартно и одновременно взаимодействовать с интерактивными элементами (кликерами) в смартфонах, синхронное ядро СУБД Oracle любого банка моментально лопнет от перегрузки строк и упадет в каскадный дедлок `ORA-00060`. «Гребцы на галерах согласований» не способны решить эту проблему штатными методами.

### 💡 Архитектурное решение КиберДеда
Технология полностью исключает нагрузку на СУБД за счет паттерна **«Холодного железа»** и динамических асинхронных микрокапсул в ОЗУ:
1. **Динамическое выделение памяти:** Микрокапсула загружается в оперативную память сервера *только* тогда, когда владелец карты реально открывает в приложении детский счет. Нет детей — сервер не тратит ни одного байта.
2. **Локальный In-Memory процессинг:** Все детские баталии, клики, нанизывания бонусов птичками и движение машинок обрабатываются со скоростью процессора строго в ОЗУ. База данных банка днем отдыхает.
3. **Пакетная консолидация:** В конце дня или недели капсула делает всего один легкий зашифрованный выстрел пачкой данных весом в 192 байта в Oracle, что полностью легально для Центробанка (ЦБ РФ), но избавляет систему от терабайтов мусора в бинарных дампах.

### 🎮 Реализованные игровые контуры (sota_duel_v3.go):
* **«Живой грузовик» и «Птичья почта»:** Раздельный гейминг для мальчиков и девочек, исключающий сваливание денег в один общий сундук.
* **Программный родительский стоп-кран:** Блокировка мамой игровых бонусов на 24 часа за непослушание (с TTL-таймером автоматического дедовского прощения).
* **Семейная интерактивная дуэль «Кто быстрее»:** Игра на одном телефоне на копейки. Ограничение сессии строго в 1 минуту (пока мама сидит у парикмахера). Маша тапает и стопорит грузовик Артема на 2 секунды, получая 10 копеек, а Артем ловит её голубя. Проигравшему летит утешительные 5 копеек, чтобы не было обид.

---

## [EN] ENGLISH VERSION

### 🎯 The Core Concept ("The Rutabaga Principle")
Traditional banking micro-savings apps for adults are boring and slow. This solution shifts the automatic rounding of parents' everyday expenses (e.g., buying rutabaga or groceries from 555 to 560 rubles) into the realm of emotional children's gamification. The product builds deep brand loyalty in the younger generation from day one, securing millions of future customers through warm childhood memories.

### 🛠 The Banking IT-Infrastructure Problem
If millions of children simultaneously tap interactive clicker elements on their smartphones, the synchronous transaction core of any traditional DBMS (like Oracle) will instantly crash under row-lock overhead and slide into a catastrophic deadlock cascade (`ORA-00060`). Corporate "galley-rowers" are physically incapable of solving this bottleneck using legacy methods.

### 💡 CyberDed's Architectural Solution
This technology completely eliminates database stress using the **"Cold Hardware"** pattern and dynamic, asynchronous in-memory micro-capsules:
1. **Dynamic Memory Allocation:** The micro-capsule is spun up in the server's RAM *only* when a parent actually activates a child account in the app. No child account means zero bytes wasted.
2. **Local In-Memory Processing:** All children's battles, taps, birds threading bonuses, and moving trucks are processed at CPU speeds strictly within RAM. The bank's main database rests during the day.
3. **Batch Consolidation:** At the end of the day or week, the capsule executes a single, lightweight, encrypted 192-byte batch transaction into Oracle. This is fully compliant with Central Bank regulations while sparing the system from terabytes of junk binary memory dumps.

### 🎮 Implemented Game Loops (sota_duel_v3.go):
* **"Live Truck" & "Bird Post":** Segregated gamification for boys and girls, preventing assets from cluttering a single ledger rows.
* **Parental Break-Lever:** A feature allowing moms to freeze game bonuses for 24 hours via a TTL-timer (with automatic "grandfatherly" forgiveness).
* **Interactive Family Duel "Who is Faster":** A 1-copeck duel on a single shared screen. Session length is strictly limited to 1 minute (while mom is at the hairdresser). Masha taps to freeze Artem’s truck for 2 seconds to win 10 copecks, while Artem tries to intercept her dove. A 5-copeck consolation prize is given to the loser to avoid tears.

---
### ⚖️ Licensing & Commercial Terms (ИП Чернов И.В.)
* **License:** MIT (Open for evaluation, adaptation, and internal corporate software extension).
* **Alfa-Bank, VTB, MTS Bank:** Available for turnkey enterprise integration under fair market value through official individual entrepreneurship contract.
* **Sberbank:** Special protective tariff — **exactly 10x market price** due to legacy corporate grudges. Pay the 10x premium or continue rowing your leaky legacy boats.
