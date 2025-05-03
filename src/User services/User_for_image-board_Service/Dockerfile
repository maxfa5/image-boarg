FROM eclipse-temurin:21-jdk-jammy

# Устанавливаем рабочую директорию
WORKDIR /app
# Копируем JAR-файл в контейнер
COPY build/libs/User_for_image-board_Service-1.0-SNAPSHOT.jar app.jar

# Открываем порт, который использует Spring Boot
EXPOSE 8001
# Команда для запуска приложения
ENTRYPOINT ["java", "-jar", "app.jar"]