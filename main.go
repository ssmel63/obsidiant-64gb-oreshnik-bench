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

type JuniorAccount struct {
	ParentID     string    `json:"parent_id"`
	ChildCardID  string    `json:"child_card_id"`
	Accumulated  float64   `json:"accumulated"`
	LastActivity time.Time `json:"last_activity"`
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

func (tc *TransactionCapsule) RegisterChildClick(parentID, childCardID string, roundDelta float64) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	bucketKey := fmt.Sprintf("%s_%s", parentID, childCardID)
	account, exists := tc.ActiveRows[bucketKey]
	if !exists {
		account = &JuniorAccount{
			ParentID:    parentID,
			ChildCardID: childCardID,
		}
		tc.ActiveRows[bucketKey] = account
	}

	account.Accumulated += roundDelta
	account.LastActivity = time.Now()
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
	fmt.Println(">>> Capsule sealed successfully.")

	return cipherText, nil
}

func main() {
	fmt.Println(">>> Launching 'Chernov-Core' Transaction Capsule...")
	capsule := NewTransactionCapsule()

	capsule.RegisterChildClick("Parent_User_77", "Junior_Card_99", 2.50)
	capsule.RegisterChildClick("Parent_User_77", "Junior_Card_99", 5.10)

	securePackage, err := capsule.EncryptAndPackWeeklyData()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf(">>> Encrypted package ready. Size: %d bytes.\n", len(securePackage))
}
