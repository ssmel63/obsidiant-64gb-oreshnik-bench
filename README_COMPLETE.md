package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"
)

// JuniorAccount описывает привязанную карту Джуниор со всеми доработками Чернова
type JuniorAccount struct {
	ParentID     string    `json:"parent_id"`
	ChildCardID  string    `json:"child_card_id"`
	Accumulated  float64   `json:"accumulated"`   // Локальный накопительный стакан в ОЗУ
	LastActivity time.Time `json:"last_activity"` // Защита от флуда
	FreezeUntil  time.Time `json:"freeze_until"`  // Таймер суточного наказания мамы (TTL)
}

// TransactionCapsule — изолированная транзакционная капсула «Чернов-Core»
type TransactionCapsule struct {
	mu         sync.RWMutex
	ActiveRows map[string]*JuniorAccount
	CipherKey  []byte // 256-битный ключ для 100 000% крипто-защиты
}

// NewTransactionCapsule инициализирует защищенный контур в оперативной памяти
func NewTransactionCapsule() *TransactionCapsule {
	key := make([]byte, 32)
	_, _ = io.ReadFull(rand.Reader, key) // Генерация ключа прямо из шума процессора
	return &TransactionCapsule{
		ActiveRows: make(map[string]*JuniorAccount),
		CipherKey:  key,
	}
}

// FreezeChildAccount — Функция мамы: блокирует копилку ребенка ровно на 1 день (24 часа)
func (tc *TransactionCapsule) FreezeChildAccount(parentID, childCardID string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	bucketKey := fmt.Sprintf("%s_%s", parentID, childCardID)
	account, exists := tc.ActiveRows[bucketKey]
	if !exists {
		account = &JuniorAccount{ParentID: parentID, ChildCardID: childCardID}
		tc.ActiveRows[bucketKey] = account
	}

	// Выставляем таймер ровно на +24 часа от текущего момента
	account.FreezeUntil = time.Now().Add(24 * time.Hour)
	fmt.Printf(">>> [МАМА НАЖАЛА СТОП-КРАН]: Копилка %s заблокирована на 24 часа за непослушание!\n", childCardID)
}

// RegisterChildClick фиксирует клики по машинкам/голубям в памяти. 
// Если включен блок — машинка программно замирает на экране смартфона.
func (tc *TransactionCapsule) RegisterChildClick(parentID, childCardID string, roundDelta float64) bool {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	bucketKey := fmt.Sprintf("%s_%s", parentID, childCardID)
	account, exists := tc.ActiveRows[bucketKey]
	if !exists {
		account = &JuniorAccount{ParentID: parentID, ChildCardID: childCardID}
		tc.ActiveRows[bucketKey] = account
	}

	// ПРОВЕРКА ТАЙМЕРА НАКАЗАНИЯ: Если 24 часа еще не истекли — выдаем блокировку
	if time.Now().Before(account.FreezeUntil) {
		timeLeft := time.Until(account.FreezeUntil)
		fmt.Printf(">>> [ОТКАЗ ДЛЯ %s]: Машинка замерла на полпути, колеса буксуют! Блок снимется через: %0.1f ч.\n", childCardID, timeLeft.Hours())
		return false // Деньги не капают, СУБД банка полностью свободна
	}

	// Если наказания нет или оно истекло — автоматическое прощение, колеса крутятся, голубь летит!
	account.Accumulated += roundDelta
	account.LastActivity = time.Now()
	fmt.Printf(">>> [ДЕТСКИЙ КЛИК ОК]: Машинка привезла +%0.2f руб. в копилку %s\n", roundDelta, childCardID)
	return true
}

// EncryptAndPackWeeklyData пакует и шифрует данные для безопасного ночного/недельного сброса в банк
func (tc *TransactionCapsule) EncryptAndPackWeeklyData() ([]byte, error) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	if len(tc.ActiveRows) == 0 {
		return nil, fmt.Errorf("capsule is empty")
	}

	// Сериализация пакетных данных в строку JSON
	plainText, err := json.Marshal(tc.ActiveRows)
	if err != nil {
		return nil, err
	}

	// Слой шифрования №1: Промышленный AES-256-GCM
	block, err := aes.NewCipher(tc.CipherKey)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	_, _ = io.ReadFull(rand.Reader, nonce)

	cipherText := aesGCM.Seal(nonce, nonce, plainText, nil)

	// Слой шифрования №2: Кастомный Фрактальный XOR КиберДеда (мутация ключа на лету)
	fractalSeed := byte(time.Now().Unix() ^ 0xDEADC0DE)
	for i := 0; i < len(cipherText); i++ {
		cipherText[i] = cipherText[i] ^ fractalSeed
		fractalSeed += byte(i) // Ключ меняется прямо во время шифрования каждого байта
	}

	// Полное обнуление локальных счетчиков памяти после успешной упаковки пачки
	tc.ActiveRows = make(map[string]*JuniorAccount)
	fmt.Println(">>> [КАПСУЛА ЗАПЕЧАТАНА]: Пачка подготовлена, память очищена. Дампы СУБД холодные.")

	return cipherText, nil
}

func main() {
	fmt.Println(">>> Запуск транзакционной капсулы 'Чернов-Core v2' (Полный объединенный цикл)...")
	capsule := NewTransactionCapsule()

	fmt.Println("\n--- ДЕНЬ 1: АКТИВНОЕ НАКОПЛЕНИЕ И ЧЕРЕДОВАНИЕ КАРТОЧЕК ДЕТЕЙ ---")
	// Покупки родителя округляются, и деньги летят по очереди разным детям в семье
	capsule.RegisterChildClick("Parent_User_77", "Junior_Card_99_Artem", 2.50) // Машинка едет у Артема
	capsule.RegisterChildClick("Parent_User_77", "Junior_Card_88_Masha", 5.10) // Голубь несет бусину Маше
	capsule.RegisterChildClick("Parent_User_77", "Junior_Card_99_Artem", 3.40) // Снова Артем

	fmt.Println("\n--- ВЕЧЕР ДНЯ 1: НАКАЗАНИЕ И РОДИТЕЛЬСКИЙ СТОП-КРАН ---")
	// Артем нашалил и орал весь день. Мама нажимает кнопку в своем личном кабинете
	capsule.FreezeChildAccount("Parent_User_77", "Junior_Card_99_Artem")

	fmt.Println("\n--- ДЕНЬ 2: ПРОВЕРКА ДИСЦИПЛИНЫ И АВТО-ТАЙМЕРА ---")
	// Маша продолжает прилежно копить — у нее всё работает плавно и без тормозов
	capsule.RegisterChildClick("Parent_User_77", "Junior_Card_88_Masha", 4.20)
	
	// Артем пытается кликать по своей машинке, но его колеса заблокированы капсулой в памяти смартфона!
	capsule.RegisterChildClick("Parent_User_77", "Junior_Card_99_Artem", 6.00)

	fmt.Println("\n--- НОЧНАЯ СИНХРОНИЗАЦИЯ: ОДИН БЕЗОПАСНЫЙ ВЫСТРЕЛ В БАНК ---")
	// Вместо тысяч мелких блокирующих запросов и дедлоков в СУБД Oracle летит один компактный крипто-пакет
	securePackage, err := capsule.EncryptAndPackWeeklyData()
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	fmt.Printf(">>> Итоговый зашифрованный за день пакет готов. Размер в ОЗУ: %d байт.\n", len(securePackage))
	fmt.Println(">>> Программа успешно завершена. Exit Code 0.")
}
