Файлы для миграции через golang-migrate

По умолчанию решил использовать функционал библиотеки Gorm

НО в случае необходимости, код для докера:

```
version: '3'

services:
  db:
    build: ./db
    env_file:
      - .env
    ports:
      - "4455:4455"

  migrate:
    image: golang-migrate/migrate
    depends_on:
      - db
    env_file:
      - .env
    volumes:
      - ./migrations:/migrations
    command: -source file:///migrations -database postgres://user:password@db:4455/database up

  app:
    build: .
    depends_on:
      - db
      - migrate
    env_file:
      - .env
    ports:
      - "8093:8093"
    command:
      - sh
      - -c
      - "sleep 3 && ./main"
      ```