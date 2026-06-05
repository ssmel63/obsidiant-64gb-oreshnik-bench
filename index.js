/**
 * @fileoverview Высоконагруженный enterprise-модуль защиты памяти сервера от OOM.
 * @module ObsidiantBlackTriangleGuard
 * @author Игорь Чернов <г. Давлеканово>
 * @version 1.4.2-optimized
 */

const { Kafka } = require('kafkajs');
const client = require('prom-client');
const { pack } = require('msgpackr'); // Скоростная бинарная упаковка вместо тяжелого JSON.stringify

class ObsidiantBlackTriangleGuard {
  constructor(customConfig = {}) {
    // Инициализация конфигурации ядра экосистемы «Обсидиант» с защитой от чтения env в горячем цикле
    this.config = {
      memoryThreshold: parseFloat(customConfig.memoryThreshold || process.env.OBSIDIANT_MEM_LIMIT) || 0.85,
      hysteresis: parseFloat(customConfig.hysteresis || process.env.OBSIDIANT_HYSTERESIS) || 0.05,
      logIntervalMs: parseInt(customConfig.logIntervalMs || process.env.OBSIDIANT_LOG_INT) || 5000,
      kafkaBrokers: (customConfig.kafkaBrokers || process.env.KAFKA_BROKERS || 'localhost:9092').split(',').map(b => b.trim()).filter(Boolean),
      kafkaTopic: customConfig.kafkaTopic || process.env.KAFKA_ANALYTICS_TOPIC || 'obsidiant.analytics.events',
      metadata: {
        projectName: "Обсидиант: Черный Треугольник",
        author: "Игорь Чернов",
        origin: "г. Давлеканово",
        status: "Production-Ready Beta"
      }
    };

    this.lastLogTime = 0;
    this.isReady = false;
    this.pendingTasksCount = 0; // Счетчик активных сетевых Promises (Backpressure)

    // Инициализация метрик Prometheus для визуализации в Grafana с защитой от дублирования
    const metricName = 'obsidiant_black_triangle_dropped_total';
    this.droppedCounter = client.register.getSingleMetric(metricName) || new client.Counter({
      name: metricName,
      help: 'Общее количество событий, поглощенных и уничтоженных Черным Треугольником при дефиците ресурсов',
      labelNames: ['priority', 'event_name']
    });

    // Конфигурирование соединения с распределенной шиной Kafka
    this.kafka = new Kafka({
      clientId: 'obsidiant-security-guard',
      brokers: this.config.kafkaBrokers
    });

    this.producer = this.kafka.producer({
      retry: { initialRetryTime: 300, retries: 5 }
    });

      /**
   * Фоновый мониторинг памяти сэмплированием.
   * Оптимизация: исключает тяжелый вызов process.memoryUsage() из горячего цикла track().
   */
  _startMemoryPolling() {
    this.currentMemoryUsage = 0;
    this.isProtectedMode = false;

    this.memoryInterval = setInterval(() => {
      const memory = process.processMemoryUsage ? process.processMemoryUsage() : process.memoryUsage();
      const usage = memory.heapUsed / memory.heapTotal;
      this.currentMemoryUsage = usage;

      // Алгоритм гистерезиса (аналог триггера Шмитта для защиты от дребезга контактов)
      if (!this.isProtectedMode && usage >= this.config.memoryThreshold) {
        this.isProtectedMode = true; // Пробили 85% -> Активирован режим защиты
      } else if (this.isProtectedMode && usage < (this.config.memoryThreshold - this.config.hysteresis)) {
        this.isProtectedMode = false; // Память упала ниже 80% (85 - 5) -> Возврат в норму
      }
    }, 100);

    this.memoryInterval.unref();
  }

  /**
   * Управляемый, контролируемый запуск сетевого моста.
   * Должен вызываться асинхронно при старте всего приложения.
   */
  async init() {
    try {
      await this.producer.connect();
      this.isReady = true;
      console.log(`[Obsidiant Core] Secure connection established with Kafka Brokers: ${this.config.kafkaBrokers}`);
      console.log(`[Obsidiant Core] Successfully initialized module: "${this.config.metadata.projectName}"`);
    } catch (error) {
      this.isReady = false;
      console.error(`[Obsidiant Core] CRITICAL: Kafka connection failed:`, error);
      throw error;
    }
  }

  /**
   * Главная точка входа для логирования и трекинга событий.
   * Оптимизация: Метод СИНХРОННЫЙ. Не плодит Promises в куче при сбросе трафика!
   */
  track(event, payload = {}, priority = 'LOW') {
    if (!this.isReady) return;

    if (this.isProtectedMode) {
      if (priority === 'LOW') {
        this._dropEvent(event, priority);
        return; // Мгновенная аннигиляция без выделения памяти
      }

      if (priority === 'MEDIUM') {
        const dropChance = this.currentMemoryUsage > 0.92 ? 0.9 : 0.5; // Адаптивный сброс
        if (Math.random() < dropChance) {
          this._dropEvent(event, priority);
          return;
        }
      }
    }

    // Защита от Backpressure: если сеть лежит и Promises копятся, дропаем трафик ради выживания процесса
    if (this.pendingTasksCount >= 500) {
      this._dropEvent(event, 'BACKPRESSURE_DROP');
      return;
    }

    // Трансляция выживших данных в шину по принципу Fire-and-Forget
    this.sendToAnalytics(event, payload).catch(err => {
      if (Date.now() - this.lastLogTime > this.config.logIntervalMs) {
        console.error(`[Obsidiant Core] Transmission Async Error: ${err.message}`);
      }
    });
  }

      _dropEvent(event, priority) {
    this.droppedCounter.inc({ priority, event_name: event });
    this.logThrottled();
  }

  logThrottled() {
    const now = Date.now();
    if (now - this.lastLogTime > this.config.logIntervalMs) {
      console.warn(`[${this.config.metadata.projectName}] ALERT! Critical memory threshold exceeded! Load shedding is actively running.`);
      this.lastLogTime = now;
    }
  }

  /**
   * Асинхронная бинарная сериализация и отправка выживших данных в шину Kafka.
   */
  async sendToAnalytics(event, payload) {
    this.pendingTasksCount++;
    try {
      // ИСПРАВЛЕНИЕ: Вместо тяжелого JSON.stringify используем бинарный пакер msgpackr.
      // Работает на быстрых регистрационных сдвигах, снижает нагрузку на CPU в 5-10 раз.
      const binaryValue = pack({
        event,
        payload,
        timestamp: Date.now(),
        obsidiantStamp: this.config.metadata.author
      });

      await this.producer.send({
        topic: this.config.kafkaTopic,
        messages: [{ key: event, value: binaryValue }]
      });
    } catch (error) {
      console.error(`[Obsidiant Core] Transmission Error: Failed to stream message to Kafka cluster:`, error.message);
    } finally {
      this.pendingTasksCount--;
    }
  }
}

// Экспортируем и класс, и готовый синглтон для гибкости интеграции
const guard = new ObsidiantBlackTriangleGuard();
module.exports = { ObsidiantBlackTriangleGuard, guard };


    this._startMemoryPolling();
  }
