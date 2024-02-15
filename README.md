# 🔐 Goph Keeper

## 💡 Overview
Goph Keeper is app for storing password, cards, texts and files. You can store locally and sync data between different clients whenever you want.

## 💻 Technologies

- **Language:** Go
- **Database:** Postgres
- **Documentation:** godoc, Swagger 2.0

## ▶️ Getting started

To get started with the Goph Keeper, follow these steps:

1. **Clone the Repository:**
```shell
git clone https://github.com/MowlCoder/goph-keeper.git
```
2. **Install Dependencies:**
```shell
go get .
```
3. **Configure Settings:** Create an `.env.client` and `.env.server` files and populate them based on the `.env.client.example` and `.env.server.example` files
4. **Run application:**
```shell
go run ./cmd/server/main.go
```
```shell
go run ./cmd/client/main.go
```

## 📝 Documentation

Documentation is available in the [docs](/docs) directory.