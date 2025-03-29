import org.gradle.internal.impldep.org.junit.experimental.categories.Categories.CategoryFilter.exclude

plugins {
    java
    id("org.springframework.boot") version "3.4.3"
    id("io.spring.dependency-management") version "1.1.7"
    id("org.hibernate.orm") version "6.4.4.Final"
}

val springBootVersion = "3.4.3"

group = "org.image_board"
version = "1.0-SNAPSHOT"
repositories {
    mavenCentral()
}



dependencies {
    implementation("org.springframework.boot:spring-boot-starter-web")
    implementation("org.springframework.boot:spring-boot-starter-data-jpa")
    implementation ("org.springframework.security:spring-security-crypto:6.4.4")
    implementation ("org.springframework.boot:spring-boot-starter-security")
//    implementation ("org.springframework.boot:spring-boot-starter-rsocket")
    implementation ("io.jsonwebtoken:jjwt-api:0.12.6")
    runtimeOnly ("io.jsonwebtoken:jjwt-impl:0.12.6")
    runtimeOnly ("io.jsonwebtoken:jjwt-jackson:0.12.6")
    implementation("org.springframework.boot:spring-boot-starter-thymeleaf") // Для фронтенда
    // Hibernate (уже включен в spring-boot-starter-data-jpa, но можно указать версию)
    implementation("org.hibernate.orm:hibernate-core")
    implementation("jakarta.persistence:jakarta.persistence-api")
    // Базы данных
    runtimeOnly("org.postgresql:postgresql")
    runtimeOnly("com.h2database:h2") // Для разработки
//  Lombok
    compileOnly("org.projectlombok:lombok:1.18.34")
    annotationProcessor("org.projectlombok:lombok:1.18.34")
    testCompileOnly("org.projectlombok:lombok:1.18.34")
    testAnnotationProcessor("org.projectlombok:lombok:1.18.34")

    // Dev tools
    //  developmentOnly("org.springframework.boot:spring-boot-devtools")

    // Тестирование
    testImplementation("org.springframework.boot:spring-boot-starter-test") {
        exclude(group = "org.junit.vintage", module = "junit-vintage-engine")
    }
    testImplementation("org.junit.jupiter:junit-jupiter-api")
    testRuntimeOnly("org.junit.jupiter:junit-jupiter-engine")

    // Для фронтенда
    implementation("org.webjars:bootstrap:5.3.0")
    implementation("org.webjars:jquery:3.6.0")
}
//testImplementation("org.junit.jupiter:junit-jupiter")


//
//tasks.withType<Test> {
//useJUnitPlatform()
//}