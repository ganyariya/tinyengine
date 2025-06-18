.PHONY: test lint build clean run

# テスト実行
test:
	go test ./... -v -race -coverprofile=coverage.out

# リント実行
lint:
	golangci-lint run

# ビルド
build:
	go build -o bin/tinyengine cmd/tinyengine/main.go

# 実行
run:
	go run cmd/tinyengine/main.go

# クリーンアップ
clean:
	rm -rf bin/ coverage.out

# 依存関係の取得
deps:
	go mod tidy
	go mod download

# カバレッジレポート表示
coverage: test
	go tool cover -html=coverage.out