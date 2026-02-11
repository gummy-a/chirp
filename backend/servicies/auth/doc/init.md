## 開発手順

### 1. go modを揃える

git cloneした初回限定。
```go
go mod download
```

#### 1.1 環境変数をセットする

```bash
export JWT_SECRET_KEY=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 128 | head -n 1)
export DATABASE_URL="postgres://postgres:password@localhost:5432/postgres?sslmode=disable"
export APP_ENV="dev"
```

### 2. DBをmigrateする

#### 2.1 localにdockerを立てる

```
docker run --name postgres -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres
```

#### 2.2 golang-migrateの入手と配置(option)

- auth serviceではgolang-migrateを使っている

```bash
curl -L https://github.com/golang-migrate/migrate/releases/latest/download/migrate.linux-amd64.tar.gz \
  | tar xvz
sudo mv migrate /usr/local/bin/
```

#### 2.3 golang-migrateでmigrateする

```bash
export DATABASE_URL="postgres://postgres:password@localhost:5432/postgres?sslmode=disable"
migrate -path sql/migrations/ -database "$DATABASE_URL" up
```

#### 2.4 schemaファイルを生成する

```bash
docker run --rm --network host postgres:18.1 pg_dump --schema-only --no-owner --no-privileges --no-comments "$DATABASE_URL" > sql/schema.sql
```

### 3. SQLからgoコード生成する

#### 3.1 sqlcをインストール(option)
```go
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

#### 3.2 コード生成

```bash
cd sql
sqlc generate
```

### 4. auth.protoからコードを生成する

#### 4.1 protocol buffer コンパイラをインストール(option)

```bash
go install github.com/golang/protobuf/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

#### 4.2 コード生成

```bash
cd backend/proto/auth
protoc --go_out=../../servicies/auth/internal/adapter/grpc/ --go-grpc_out=../../servicies/auth/internal/adapter/grpc/ auth.proto
```

### 5. OpenAPIからコードを生成する

プロジェクトrootで実施する。
```go
az@debian:~/code/chirp$ go run api/generate.go
```

#### 6 新しく依存関係ができるのでインストールする

```go
go mod tidy
```
