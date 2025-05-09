# Markblog

[![xc compatible](https://xcfile.dev/badge.svg)](https://xcfile.dev)

Для поднятия локально требуется установить [encore](https://encore.dev) и запустить проект с помощью ```encore run```

## Стэк

- Encore
- Go
- PostgreSQL
- sqlc
- Vue
- Typescript
- Pinia
- TailwindCSS
- DaisyUI
- Docker

## Tasks

### run
```bash
encore run
```

### dev

```bash
sqlc generate
encore gen client --output=./webapp/frontend/src/client.ts -x=api --excluded-tags=noclient
npm run format --prefix webapp/frontend
npm run build --prefix webapp/frontend
encore run
```

### devfrontend

```bash
npm run dev --prefix webapp/frontend --watch
```
