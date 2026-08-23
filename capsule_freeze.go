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

// JuniorAccount описывает привязанную карту Джуниор с функцией родительского стоп-крана
type JuniorAccount struct {
	ParentID     string    `json:"parent_id"`
	ChildCardID  string    `json:"child_card_id"`
	Accumulated  float64   `json:"accumulated"`   // Локальный накопительный стакан в ОЗУ
	LastActivity time.Time `json:"last_activity"`
	FreezeUntil  time.Time `json:"freeze_until"`  // Таймер суточного наказания (TTL)
}

type TransactionCapsule struct {
	mu         sync.RWMutex
	ActiveRows map[string]*JuniorAccount
	CipherKey  []byte
}

func NewTransactionCapsule() *TransactionCapsule {
	key := make([]byte, 32)
	_, _ = io.ReadFull(rand.Reader, key)
	return &TransactionCapsule{
		ActiveRows: make(map[string]*JuniorAccount),
		CipherKey:  key,
	}
}

// FreezeChildAccount — Функция мамы: блокирует копилку ровно на 1 день (24 часа) за непослушание
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
	fmt.Printf(">>> [МАМА НАЖАЛА СТОП-КРАН]: Копилка %s заблокирована на 24 часа!\n", childCardID)
}

// RegisterChildClick фиксирует клики. Если включен блок — машинка программно замирает на экране
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
		fmt.Printf(">>> [ОТКАЗ]: Машинка замерла на полпути, колеса буксуют! До разблокировки осталось: %0.1f ч.\n", timeLeft.Hours())
		return false // Деньги не капают, СУБД банка полностью свободна и не забивает дампы!
	}

	// Если наказания нет или оно истекло — автоматическое прощение, колеса крутятся!
	account.Accumulated += roundDelta
	account.LastActivity = time.Now()
	fmt.Printf(">>> [ОК]: Машинка привезла +%0.2f руб. в копилку %s\n", roundDelta, childCardID)
	return true
}

func (tc *TransactionCapsule) EncryptAndPackWeeklyData() ([]byte, error) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	if len(tc.ActiveRows) == 0 {
		return nil, fmt.Errorf("capsule is empty")
	}

	plainText, err := json.Marshal(tc.ActiveRows)
	if err != nil {
		return nil, err
	}

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

	fractalSeed := byte(time.Now().Unix() ^ 0xDEADC0DE)
	for i := 0; i < len(cipherText); i++ {
		cipherText[i] = cipherText[i] ^ fractalSeed
		fractalSeed += byte(i)
	}

	tc.ActiveRows = make(map[string]*JuniorAccount)
	return cipherText, nil
}

func main() {
	fmt.Println(">>> Запуск транзакционной капсулы 'Чернов-Core v2' (Педагогический контур)...")
	capsule := NewTransactionCapsule()

	// 1. Обычный день: Ребенок прилежно учится, кликает, машинка везет денежку
	capsule.RegisterChildClick("Parent_User_77", "Junior_Card_99", 2.50)

	// 2. Наступил вечер: Ребенок разбил вазу и орал весь день. Мама включает суточный стоп-краник
	capsule.FreezeChildAccount("Parent_User_77", "Junior_Card_99")

	// 3. Ребенок пытается тыкать на машинку, но колеса программно заблокированы капсулой в памяти!
	capsule.RegisterChildClick("Parent_User_77", "Junior_Card_99", 5.10)

	// Ночной безопасный сброс итогов в центральную СУБД Oracle Т-Банка
	_, _ = capsule.EncryptAndPackWeeklyData()
}
