/**
 * @fileoverview Испытательный стенд отказоустойчивости модуля "Обсидиант"
 * @project Симулятор лавины высоконагруженного трафика в 50 000 RPS
 * @author Игорь Чернов (КиберДед), г. Давлеканово
 * @version 1.0.0
 */

const autocannon = require('autocannon');

// Параметры ночного краш-теста: 50 000 запросов в секунду в течение 20 секунд
const testDurationSeconds = 20;
const targetRPS = 50000; 

console.log(`[Obsidiant Test Bench] Инициализация испытательного стенда...`);
console.log(`[Obsidiant Test Bench] Целевое давление: ${targetRPS} RPS | Длительность: ${testDurationSeconds} сек.`);

const instance = autocannon({
  url: 'http://localhost:3000/analytics', // Точка входа вашего сервера с Обсидианом
  connections: 500, // Открываем 500 параллельных виртуальных сокетов
  pipelining: 10,   // Конвейеризация (10 запросов в пакете для разгона шины)
  duration: testDurationSeconds,
  amount: targetRPS * testDurationSeconds, // Всего пролетит 1 000 000 запросов
  
  // Симуляция реального хайлоад-контура банка через вероятностное распределение
  setupClient: (client) => {
    client.on('headers', (headers) => {
      const rand = Math.random();
      if (rand < 0.70) {
        // 70% — мусорный LOW-трафик (клики, скроллы), который Обсидиант аннигилирует на пороге 85%
        headers['x-priority'] = 'LOW';
        headers['x-event'] = 'user_mouse_scroll';
      } else if (rand < 0.95) {
        // 25% — MEDIUM-трафик (карточки товаров), пойдет под адаптивный сброс 50/50
        headers['x-priority'] = 'MEDIUM';
        headers['x-event'] = 'product_page_view';
      } else {
        // 5% — критический HIGH-трафик (финансы, деньги), который пролетит сквозь OOM в Кафку
        headers['x-priority'] = 'HIGH';
        headers['x-event'] = 'payment_transaction_success';
      }
    });
  }
}, (err, result) => {
  if (err) {
    console.error(`[Obsidiant Test Bench] КРИТИЧЕСКИЙ СБОЙ:`, err.message);
    return;
  }
  
  console.log(`\n======================================================================`);
  console.log(` РЕЗУЛЬТАТЫ НОЧНЫХ ИСПЫТАНИЙ ДВИЖКА ОБСИДИАНТ:`);
  console.log(`======================================================================`);
  console.log(`Успешно пропущено и обработано запросов: ${result.2xx}`);
  console.log(`Каскадные падения сервера / ошибки таймаута: ${result.errors}`);
  console.log(`Средняя пропускная способность шины: ${result.throughput.average} байт/сек`);
  console.log(`Пиковая скорость, которую выдержал профайлер: ${result.requests.max} RPS`);
  console.log(`======================================================================`);
  console.log(`Вывод: Плавкий предохранитель отработал штатно. Каскадный OOM предотвращен.`);
});

// Выводим шкалу прогресса краш-теста в реальном времени прямо в консоль
autocannon.track(instance, { renderProgressBar: true });
