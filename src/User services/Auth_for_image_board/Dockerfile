FROM eclipse-temurin:21-jdk-jammy

# Устанавливаем рабочую директорию
WORKDIR /app
# Копируем JAR-файл в контейнер
COPY build/libs/Auth_for_image_board-1.0-SNAPSHOT.jar app.jar

# Открываем порт, который использует Spring Boot
EXPOSE 8002
# Команда для запуска приложения
ENTRYPOINT ["java", "-jar", "app.jar"]