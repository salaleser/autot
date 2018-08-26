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
Для работы нужен файл настроек `config` со следующими ключами:
```
ver=1.0.0
templateDateFormat=2006-01-02_15-04
server=\\server-path
service=service-name
pathTemp=temp\
pathData=\\server-path\folder-name\
pathBackup=\\server-path\folder-name\backup\
pathKmis=\\server2-path\folder-name
pathSigned=\\server2-path\folder-name\folder2-name\
countdown=6
cooldown=500
```
