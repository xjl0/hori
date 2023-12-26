# Используем базовый образ golang для сборки приложения
FROM golang:1.21.5 as builder
WORKDIR /backend-app-files
# Копируем исходный код
COPY . .
# Собираем статически линкованный исполняемый файл
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./hori-bot ./main.go

# Используем базовый образ alpine для минимального образа
FROM alpine:3.18.5
# Добавляем пользователя `app`
RUN addgroup -g 1001 app && \
    adduser -u 1001 -D -G app app /home/app
# Устанавливаем ca-certificates и tzdata
RUN apk add --no-cache ca-certificates tzdata

# Устанавливаем переменную окружения для часового пояса
ENV TZ=Europe/Moscow

# Копируем собранный файл из предыдущего шага
COPY --chown=1001:1001 --from=builder /backend-app-files/hori-bot /hori-bot

# Указываем пользователя
USER app

# Устанавливаем переменные окружения для приложения
ENV GIN_MODE=release \
    BACKEND_PORT=4005

# Задаем точку входа
ENTRYPOINT ["/hori-bot"]