
# Discord Бот Hori

- Интегрирован ChatGPT 4o
- Получает текущие праздники с [сайта](https://kakoysegodnyaprazdnik.ru)

## Настройка бота

Для работы бота необходимо [создать приложение](https://discord.com/developers/applications) на сайте discord

Настройка env
```env
TI_DASHBOARD_CHANNEL=Канал куда отправляются праздники
TI_DISCORD_BOT_TOKEN=Токен дискод
TI_PROXY=url прокси для ChatGPT если сервер находится в стране где не работает сервис
TI_OPENAI_API_KEY=токен openAI
TI_SIZE_CONTEXT=количество запоминаемых сообщений в канале
```
## Команды

Бот по умолчанию запонимает последние 50 сообщений в канале.

Чтобы вызвать бота необходимо его тегнуть.