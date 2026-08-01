package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

type ChaosPacket struct {
	Type uint8 // 1 - Данные, 2 - Обломки, 3 - Мусор
}

func main() {
	fmt.Println("==================================================")
	fmt.Println("   ЗАПУСК ВЫСОКОСКОРОСТНОГО СТЕНДА «ОРЕШНИК»")
	fmt.Println("   Имитация сверхнагрузки на пределе процессора")
	fmt.Println("==================================================")
	fmt.Println("Для остановки теста нажмите Ctrl+C\n")

	var countData uint64
	var countWreck uint64
	var countGarbage uint64

	packetChannel := make(chan ChaosPacket, 100000)

	go func() {
		var i uint64
		for {
			var pType uint8
			switch i % 3 {
			case 0:
				pType = 1
			case 1:
				pType = 2
			case 2:
				pType = 3
			}
			packetChannel <- ChaosPacket{Type: pType}
			i++
		}
	}()

	go func() {
		for packet := range packetChannel {
			switch packet.Type {
			case 1:
				atomic.AddUint64(&countData, 1)
			case 2:
				atomic.AddUint64(&countWreck, 1)
			case 3:
				atomic.AddUint64(&countGarbage, 1)
			}
		}
	}()

	for {
		time.Sleep(1 * time.Second)
		
		data := atomic.SwapUint64(&countData, 0)
		wreck := atomic.SwapUint64(&countWreck, 0)
		garbage := atomic.SwapUint64(&countGarbage, 0)
		total := data + wreck + garbage

		fmt.Printf("[%s] Залп «Орешника»: ОБЩАЯ СКОРОСТЬ = %d пак/сек\n", time.Now().Format("15:04:05"), total)
		fmt.Printf("           ├── Полезные данные: %d\n", data)
		fmt.Printf("           ├── Обломки структуры: %d\n", wreck)
		fmt.Printf("           └── Отсеченный мусор:  %d\n\n", garbage)
	}
}
