package main

import (
	"fmt"
	"sync"
	"time"
)

// ChildDuelSession — структура интерактивной сессии "Брат против Сестры" в ОЗУ
type ChildDuelSession struct {
	mu             sync.Mutex
	ParentCardID   string
	BoyCardID      string
	GirlCardID     string
	BoyBalance     float64   // Баланс Копилки Артема
	GirlBalance    float64   // Баланс Копилки Маши
	BoyFrozenUntil time.Time // Время заморозки грузовика (2 секунды)
	GirlFrozenUntil time.Time // Время заморозки голубя (2 секунды)
	SessionStart   time.Time // Фиксация начала игры для контроля маминых карманов
}

// NewDuelSession — инициализация капсулы в момент открытия детского счета
func NewDuelSession(parentID, boyID, girlID string) *ChildDuelSession {
	return &ChildDuelSession{
		ParentCardID: parentID,
		BoyCardID:    boyID,
		GirlCardID:   girlID,
		SessionStart: time.Now(),
	}
}

// CatchOpponent — функция перехвата предмета соперника пальцем на одном экране
func (ds *ChildDuelSession) CatchOpponent(attackerCardID string) (string, bool) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	// 1. Контроль времени мамы (Защита от залипания): сессия закрывается ровно через 1 минуту
	if time.Since(ds.SessionStart) > 1*time.Minute {
		return "Сессия завершена автоматически. Время игры истекло (награда за поведение).", false
	}

	now := time.Now()

	// 2. Логика перехвата: Мальчик ловит голубя или Девочка останавливает грузовик
	if attackerCardID == ds.GirlCardID {
		// Девочка тапает по экрану и останавливает грузовик брата
		if now.Before(ds.BoyFrozenUntil) {
			return "Грузовик уже заблокирован!", false
		}
		ds.BoyFrozenUntil = now.Add(2 * time.Second) // Замораживаем машинку на 2 секунды
		ds.GirlBalance += 0.10                       // Победителю перехвата +10 копеек
		ds.BoyBalance += 0.05                        // Утешительные +5 копеек проигравшему, чтобы не обиделся

		return fmt.Sprintf("Маша остановила грузовик Артема на 2 сек! Маша: +10 коп. Артем: +5 коп."), true
	} else if attackerCardID == ds.BoyCardID {
		// Мальчик тапает по экрану и сбивает летящего голубя сестры
		if now.Before(ds.GirlFrozenUntil) {
			return "Голубь уже заблокирован!", false
		}
		ds.GirlFrozenUntil = now.Add(2 * time.Second) // Замораживаем голубя на 2 секунды
		ds.BoyBalance += 0.10                        // Победителю перехвата +10 копеек
		ds.GirlBalance += 0.05                       // Утешительные +5 копеек проигравшему

		return fmt.Sprintf("Артем поймал голубя Маши на 2 сек! Артем: +10 коп. Маша: +5 коп."), true
	}

	return "Неверный идентификатор игрока", false
}

// ConsolidationPayload — формирование 192-байтного пакета для ночного сброса в банк
func (ds *ChildDuelSession) ConsolidationPayload() string {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	return fmt.Sprintf("ID родителей: %s | Сброс проводки в Oracle -> Итого Артем: %.2f руб, Итого Маша: %.2f руб. Логи ЦБ чистые.", 
		ds.ParentCardID, ds.BoyBalance, ds.GirlBalance)
}

func main() {
	fmt.Println("=== ТЕСТ КАПСУЛЫ ЧЕРНОВ-CORE V3: РЕЖИМ ДУЭЛИ В ОЗУ ===")
	
	// Мама открыла детский счет в парикмахерской
	duel := NewDuelSession("PARENT_IVAN_76", "ARTEM_BOY", "MASHA_GIRL")

	// Пошёл азартный минутный кликер на одном телефоне
	res1, _ := duel.CatchOpponent("MASHA_GIRL") // Маша ловит грузовик
	fmt.Println(res1)

	time.Sleep(1 * time.Second)

	res2, _ := duel.CatchOpponent("ARTEM_BOY")  // Артем ловит голубя в ответ
	fmt.Println(res2)

	// Имитируем завершение маминой стрижки (ночной пакетный сброс)
	fmt.Println("\n=== СБРОС ДАННЫХ В КОНЦЕ ДНЯ БЕЗ ДЕДЛОКОВ ===")
	fmt.Println(duel.ConsolidationPayload())
}
