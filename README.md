
# Discord Бот Hori

- Автоматически приветствует
- Ставит реакцию на ссылки youtube и coub для рейтинга
- Получает текущие праздники с [сайта](https://kakoysegodnyaprazdnik.ru)
- Получает последнюю новость с [сайта](https://shikimori.one)
- Вычисляет время просмотра серий
- Автоматически переводит сообщения на русский язык
- Отправляет в чат ссылку на новое созданное мероприятие (событие) сервера
## Настройка бота

Для работы бота необходимо [создать приложение](https://discord.com/developers/applications) на сайте discord

Настройка env
```env
DGU_TOKEN=Токен приложения
DISCORD_CODEINVITE=Код приглашения на сервер
DISCORD_HELLOEMOTE=ID эмоции сервера
DISCORD_LASTED=ID последней новости Шики
DISCORD_MAINCHANNEL=ID основного чата
DISCORD_MEDIACHANNEL=ID медиа чата
DISCORD_NEWSCHANNEL=ID чата новостей
DISCORD_TESTCHANNEL=ID чата для тестов
```
## Команды

- календарь
- Сколько по времени 1 серия? (Сколько по времени 12 серий?)
- бот жив?