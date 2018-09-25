# Бот для slack для автоматизации кое-чего для кое-кого.

## Серверная версия
Это версия бота для запуска на сервере.
Если в качестве аргумента исполняемому файлу будет указан правильный токен бота
вида  `xoxb-000000000000-000000000000-AAAAAAAAAAAAAAAAAAAAAAAA`, то бот
будет запущен в полноценном режиме и сможет обрабатывать команды в слэке.

## Клиентская версия
Это версия бота для запуска на клиенте.
Также висит в трее винды как и полная версия, но без слэк-бота.

## Необходимые файлы
Для работы нужен файл настроек `autot.cfg` со следующими ключами:
```
server=SERVER_NAME
service=SERVICE_NAME
data-dir=DATA_DIR
src-dir=SOURCE_DIR
dest-dir=DEST_DIR
countdown=15
cooldown=500
stopped-sound=stopped.wav
start-pending-sound=start_pending.wav
stop-pending-sound=stop_pending.wav
started-sound=started.wav
stop-poll-sound=stop-poll.wav
beep-sound=beep.wav
```
