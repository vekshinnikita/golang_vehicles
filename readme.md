## Getting started

1. Добавить в корневую директорию файл `.env` со структурой из файла `.env.example`

2. Запустить бд в докер командой 
```sh
docker run --name=vehicle-db -e POSTGRES_PASSWORD='admin' -v "$(pwd)/pgdata":/var/lib/postgresql/data -p 5432:5432 -d --rm  postgres
```

3. Применить миграции командой
```sh
migrate -path ./schema -database 'postgres://postgres:admin@localhost:5432/postgres?sslmode=disable' up
```

4. Запустить проект командой 
```sh 
air
```