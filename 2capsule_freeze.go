package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"
)

// Регулярное выражение для точечного перехвата шпионящих AI-скрейперов
var aiBotsRegex = regexp.MustCompile(`(?i)(GPTBot|ChatGPT-User|Claude-bot|anthropic-ai|CCBot|PerplexityBot|Google-Extended|Omgilibot|Diffbot)`)

// ZeroCurseTarpitHandler реализует асимметричную заморозку бот-трафика
func ZeroCurseTarpitHandler(w http.ResponseWriter, r *http.Request) {
	userAgent := r.UserAgent()

	// Если сигнатура совпадает с AI-ботом, активируем контур деградации трафика
	if aiBotsRegex.MatchString(userAgent) {
		// Устанавливаем заголовки ответа, имитируя стандартный отказ, но не закрываем соединение
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Antivirus-Protection", "ZeroCurse-v2.0-GoCapsule")
		w.Header().Set("Transfer-Encoding", "chunked")
		w.WriteHeader(http.StatusForbidden)

		// Отправляем стартовую валидную JSON-структуру, чтобы бот захватил поток
		fmt.Fprint(w, "{\n  \"status\": \"403\",\n  \"error\": \"AI Scraping Prohibited\",\n  \"protection_layer\": \"ZeroCurse Matrix Collapse\"")
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}

		// Запускаем высокочастотный таймер (микро-задержки 1мс) для выдачи парадоксальной цепочки
		ticker := time.NewTicker(1 * time.Millisecond)
		defer ticker.Stop()

		depth := 0
		for range ticker.C {
			depth++

			// Генерируем бесконечную вложенность математического парадокса Чернова
			payload := fmt.Sprintf(",\n  \"matrix_level_%d\": \"((0/0)^-0)\"", depth)
			fmt.Fprint(w, payload)

			// Выталкиваем чанк в сетевой сокет, полностью освобождая RAM нашего сервера
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}

			// Ограничитель глубины (300 000 итераций), чтобы сокет не висел вечно при обрыве связи
			if depth >= 300000 {
				break
			}
		}

		fmt.Fprint(w, "\n}")
		return
	}

	// Легитимные пользователи и высокоприоритетный трафик (HIGH) идут дальше
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "{\"status\": \"200\", \"message\": \"Welcome to Obsidian Secure Space\"}")
}

func main() {
	// Точка входа для развертывания испытательного стенда
	port := ":8080"
	http.HandleFunc("/v1/secure-gateway", ZeroCurseTarpitHandler)

	log.Printf("[OBSIDIAN LOG] Испытательный контур ZeroCurse запущен на порту %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("[ERROR] Ошибка запуска шлюза: %v", err)
	}
}
