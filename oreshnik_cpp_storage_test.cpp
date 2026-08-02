#include <iostream>
#include <fstream>
#include <chrono>
#include <thread>
#include <cstdint>

// Наша патентная 9-байтовая структура
#pragma pack(push, 1)
struct StreamTask {
    uint64_t id;
    uint8_t status;
};
#pragma pack(pop)

int main() {
    std::cout << "==================================================" << std::endl;
    std::cout << "   СТЕНД СТУПЕНЧАТОГО РАЗГОНА НА C++ (WINDOWS)" << std::endl;
    std::cout << "==================================================" << std::endl;

    std::ofstream file("potok_storage_cpp.db", std::ios::out | std::ios::binary | std::ios::app);
    if (!file.is_open()) {
        std::cout << "Ошибка создания файла базы!" << std::endl;
        return 1;
    }

    uint64_t totalCounter = 0;
    uint64_t writeCount = 0;
    int64_t delayMicroseconds = 500;

    auto lastReportTime = std::chrono::steady_clock::now();

    while (true) {
        totalCounter++;
        writeCount++;

        StreamTask task;
        task.id = totalCounter;
        task.status = 0b0010; // Ваша бинарная маска

        file.write(reinterpret_cast<const char*>(&task), sizeof(StreamTask));

        if (!file) {
            std::cout << "\nДиск переполнен! Испытание на C++ завершено." << std::endl;
            break;
        }

        if (totalCounter % 10000 == 0) {
            file.flush();
        }

        auto currentTime = std::chrono::steady_clock::now();
        auto elapsedTime = std::chrono::duration_cast<std::chrono::seconds>(currentTime - lastReportTime).count();

        if (elapsedTime >= 1) {
            std::cout << "Задержка: " << delayMicroseconds << " мкс | СКОРОСТЬ НА C++: " << writeCount << " задач/сек" << std::endl;
            
            writeCount = 0;
            lastReportTime = currentTime;

            if (delayMicroseconds > 100) {
                delayMicroseconds -= 100;
            } else if (delayMicroseconds > 10) {
                delayMicroseconds -= 10;
            } else if (delayMicroseconds > 0) {
                delayMicroseconds -= 1;
            }
        }

        if (delayMicroseconds > 0) {
            std::this_thread::sleep_for(std::chrono::microseconds(delayMicroseconds));
        }
    }

    file.close();
    return 0;
}
