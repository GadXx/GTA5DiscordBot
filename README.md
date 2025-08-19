## Установка и запуск

### 1. Обновляем пакеты
```bash
sudo apt update && sudo apt upgrade -y
```

### 2. Устанавливаем Docker и Docker Compose
```bash
sudo apt install -y docker.io docker-compose
```

### Добавим текущего пользователя в группу docker (чтобы не писать sudo):
```bash
sudo usermod -aG docker $USER
```

### 3. Клонируем репозиторий
```bash
git clone https://github.com/GadXx/GTA5DiscordBot.git
```
```bash
cd GTA5DiscordBot
```

### 4. Настраиваем переменные окружения
```bash
nano
```

#### Пример .env файла
```bash
BOT_TOKEN=...api токен бота
DB_PATH=/app/db/database.sqlite
GUILD_ID=...id сервера

LOG_CHANNEL_ID=...id чата в который будет писаться лог
ROLE_VAC_ID=...id роли которая выдается на время(отпуск)
ROLE_SANCTION_ID=...id роли которую может выдать только "главный"
ROLE_LEADER_ID=...id роли "главный"
DEFAULT_ROLE_ID=...id роли которая выдается по дефолту при заходе на сервер
```

### 5. Запуск бота
#### Сбилдить
```bash
docker-compose build --no-cache
```
#### Запустить
```bash
docker-compose up -d
```
#### Посмотреть логи
```bash
docker-compose logs -f
```
#### Остановить бота
```bash
docker-compose down
```