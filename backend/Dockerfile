FROM gradle:8.7-jdk17-alpine AS build
WORKDIR /home/gradle/src
COPY --chown=gradle:gradle . .
RUN gradle build --no-daemon

FROM openjdk:17-alpine
WORKDIR /app
COPY --from=build /home/gradle/src/build/libs/*.jar /app/app.jar
EXPOSE 8080
ENTRYPOINT ["java", "-jar", "/app/app.jar"]
