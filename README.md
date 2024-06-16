
# Discord Бот Hori

- Интеграция с GPT4o
- Получает текущие праздники с [сайта](https://kakoysegodnyaprazdnik.ru)

## Настройка бота

Для работы бота необходимо [создать приложение](https://discord.com/developers/applications) на сайте discord

Настройка env
```env
TI_SIZE_CONTEXT=количество запоминаемых сообщений в канале
TI_DASHBOARD_CHANNEL=канал куда отсылать сообщения о праздниках
TI_PROXY=прокси сервер если не работет chat gpt в стране нахождения сервера
TI_DISCORD_BOT_TOKEN=токен бота дискорд
TI_OPENAI_API_KEY=токен openai
```
## Команды

Чтобы вызвать бота необходимо его тегнуть