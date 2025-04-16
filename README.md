# markblog

[![xc compatible](https://xcfile.dev/badge.svg)](https://xcfile.dev)

Проект для ознакомительной практики

## Tasks

### dev

```bash
docker compose down
sqlc generate
npm run build --prefix web
docker compose up -d --build
```

### gen

```bash
sqlc generate
npm run build --prefix web
```

### run

```bash
docker compose up -d --build
```

### sqlc

```bash
sqlc generate
```
