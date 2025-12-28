# テスト実行ガイド

## 概要

本ドキュメントでは、Firebase カスタム認証 API のサーバーサイドにおける統合テストの実行方法について説明します。

## テスト方針

### 古典テスト（Classical Testing）アプローチ

本プロジェクトでは、**モックを使用しない古典的なテストアプローチ**を採用しています。

**採用理由:**

- 実際の動作環境に近い状態でテストを実行
- Firebase Emulator を使用することで、本番環境と同等の挙動を確認
- モックのメンテナンスコストを削減
- 統合的な動作検証が可能

### 使用技術

- **Firebase Emulator Suite**: Firestore および Firebase Auth のエミュレーター
- **httptest**: HTTP ハンドラーのテスト
- **Go 標準テストパッケージ**: `testing` パッケージ

## テスト環境

### 必要な環境変数

```bash
FIRESTORE_EMULATOR_HOST=localhost:8080
FIREBASE_AUTH_EMULATOR_HOST=localhost:9099
```

### 前提条件

- Firebase Emulator Suite が起動していること
- Docker コンテナで Firebase Emulator が実行されていること

### Firebase Emulator の起動確認

```bash
# Docker コンテナの状態を確認
docker ps | grep firebase

# 期待される出力例:
# custom_auth_with_firebase-firebase-1   Up X seconds
```

Emulator が起動していない場合:

```bash
# Docker Compose で起動
docker-compose up -d
```

## テスト実行方法

### 全テスト実行

```bash
# Firebase Emulator が起動していることを確認
docker ps | grep firebase

# テスト実行
FIRESTORE_EMULATOR_HOST=localhost:8080 \
FIREBASE_AUTH_EMULATOR_HOST=localhost:9099 \
go test ./... -v -cover
```

### カバレッジレポート生成

```bash
# カバレッジプロファイル生成
FIRESTORE_EMULATOR_HOST=localhost:8080 \
FIREBASE_AUTH_EMULATOR_HOST=localhost:9099 \
go test ./... -coverprofile=coverage.out

# HTML レポート生成
go tool cover -html=coverage.out -o coverage.html

# テキスト形式でカバレッジ確認
go tool cover -func=coverage.out | tail -20
```

### 特定パッケージのテスト実行

#### Repository層のみ

```bash
FIRESTORE_EMULATOR_HOST=localhost:8080 \
go test ./internal/infrastructure/persistence -v -cover
```

#### Usecase層のみ

```bash
FIRESTORE_EMULATOR_HOST=localhost:8080 \
go test ./internal/usecase -v -cover
```

#### Handler層のみ

```bash
FIRESTORE_EMULATOR_HOST=localhost:8080 \
FIREBASE_AUTH_EMULATOR_HOST=localhost:9099 \
go test ./internal/interface/handler -v -cover
```

### テストの並列実行

```bash
# 並列数を指定してテスト実行（デフォルトは CPU コア数）
FIRESTORE_EMULATOR_HOST=localhost:8080 \
FIREBASE_AUTH_EMULATOR_HOST=localhost:9099 \
go test ./... -v -cover -parallel 4
```

## トラブルシューティング

### Firebase Emulator が起動していない

**症状:**

```
Skipping integration test: FIRESTORE_EMULATOR_HOST is not set.
```

**原因:** 環境変数が設定されていないか、Firebase Emulator が起動していない

**解決方法:**

```bash
# 1. Docker コンテナを確認
docker ps | grep firebase

# 2. コンテナが起動していない場合
docker-compose up -d

# 3. 環境変数が正しく設定されているか確認
echo $FIRESTORE_EMULATOR_HOST
echo $FIREBASE_AUTH_EMULATOR_HOST
```

### 接続エラー（grpc: the client connection is closing）

**症状:**

```
rpc error: code = Canceled desc = grpc: the client connection is closing
```

**原因:** テスト実行中に Firestore クライアントが複数回クローズされている

**解決方法:**

これは cleanup 時の警告ログであり、テスト自体は成功しています。無視しても問題ありません。

ログを抑制したい場合は、`t.Cleanup()` 内のエラーログを調整してください。

### テストが遅い

**症状:** テスト実行に時間がかかる

**原因:** Firebase Emulator との通信によるオーバーヘッド

**改善方法:**

1. **並列実行を活用**
   ```bash
   go test ./... -parallel 4
   ```

2. **特定パッケージのみ実行**
   ```bash
   # 開発中は変更したパッケージのみテスト
   go test ./internal/usecase -v
   ```

3. **短いタイムアウトを設定**
   ```bash
   go test ./... -timeout 5m
   ```

**注意:** Firestore クライアントの同時接続数が多すぎるとエラーが発生する可能性があります。

### テストがスキップされる

**症状:**

```
--- SKIP: TestOTPRepository_Save (0.00s)
    otp_repository_test.go:20: Skipping integration test: FIRESTORE_EMULATOR_HOST is not set.
```

**原因:** 環境変数が設定されていない

**解決方法:**

```bash
# 環境変数をエクスポート
export FIRESTORE_EMULATOR_HOST=localhost:8080
export FIREBASE_AUTH_EMULATOR_HOST=localhost:9099

# その後テスト実行
go test ./... -v
```

### Firestore のデータが残っている

**症状:** テスト実行後に Emulator にテストデータが残っている

**原因:** `t.Cleanup()` が正しく実行されていない、または意図的に残している

**解決方法:**

```bash
# Firebase Emulator を再起動してデータをクリア
docker-compose restart firebase

# または、Emulator UI から手動で削除
# http://localhost:4000 にアクセス
```

### ポート競合エラー

**症状:**

```
Error: Port 8080 is already in use
```

**原因:** 別のプロセスが同じポートを使用している

**解決方法:**

```bash
# 1. ポートを使用しているプロセスを確認
lsof -i :8080
lsof -i :9099

# 2. プロセスを終了
kill -9 <PID>

# 3. Firebase Emulator を再起動
docker-compose restart firebase
```

## 継続的インテグレーション（CI）での実行

### GitHub Actions の例

```yaml
name: Test

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest

    services:
      firebase:
        image: firebase/emulator-suite
        ports:
          - 8080:8080
          - 9099:9099
        options: >-
          --health-cmd "curl -f http://localhost:4000 || exit 1"
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Run tests
        env:
          FIRESTORE_EMULATOR_HOST: localhost:8080
          FIREBASE_AUTH_EMULATOR_HOST: localhost:9099
        run: |
          cd server
          go test ./... -v -coverprofile=coverage.out

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./server/coverage.out
```

## ベストプラクティス

### 1. テスト実行前の確認事項

- [ ] Firebase Emulator が起動している
- [ ] 環境変数が正しく設定されている
- [ ] 前回のテストデータがクリアされている（必要に応じて）

### 2. テストの実行順序

1. **ローカル開発時**: 変更したパッケージのみテスト
2. **コミット前**: 全テストを実行
3. **PR作成時**: CI で全テストが実行されることを確認

### 3. カバレッジ目標

- 各層で **80% 以上**のカバレッジを目指す
- 新規コードは必ずテストを追加する
- カバレッジが低下する PR は避ける

### 4. テストの保守

- テストが失敗したら、まず本番コードではなくテストを疑う
- テストが脆くなってきたら、リファクタリングを検討
- テストの実行時間が長くなってきたら、並列化や最適化を検討

## 参考資料

- [Go Testing Best Practices](https://go.dev/doc/tutorial/add-a-test)
- [Firebase Emulator Suite](https://firebase.google.com/docs/emulator-suite)
- [httptest パッケージ](https://pkg.go.dev/net/http/httptest)
- [Table Driven Tests in Go](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)

---

**関連ドキュメント:**
- [テスト仕様書](./test-specifications.md) - 各層のテストケース詳細
- [セキュリティテスト](./security-testing.md) - セキュリティ対策の検証
- [設計の意思決定](./design-decisions.md) - テスト設計の背景

**最終更新日**: 2025-12-28
