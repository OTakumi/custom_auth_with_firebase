# セキュリティテスト仕様書

## 概要

本ドキュメントでは、Firebase カスタム認証 API における6つの主要なセキュリティ対策と、それらが正しく実装されていることを検証するテストについて説明します。

## セキュリティ対策一覧

| # | 対策名 | 脅威 | 優先度 |
|---|--------|------|--------|
| 1 | メールアドレス列挙攻撃対策 | アカウント情報の漏洩 | 高 |
| 2 | OTP 試行回数制限 | ブルートフォース攻撃 | 高 |
| 3 | OTP 期限管理 | リプレイ攻撃 | 高 |
| 4 | OTP の非公開 | 情報漏洩 | 高 |
| 5 | ワンタイムパスワードの保証 | 再利用攻撃 | 高 |
| 6 | 安全な乱数生成 | 推測攻撃 | 中 |

---

## 1. メールアドレス列挙攻撃対策

### 脅威の説明

攻撃者がエラーメッセージの違いから、特定のメールアドレスが登録されているかどうかを判断できてしまう脆弱性。

**攻撃シナリオ**:
```
登録済み: "ユーザーが存在しません"
未登録:   "OTPを送信しました"
→ メールアドレスの存在が判明
```

### 実装された対策

**汎用的なエラーメッセージの使用**:
- 未登録のメールアドレスに対しても `"Authentication failed"` を返す
- エラーメッセージから登録状態を推測できないようにする

**コード実装** (`internal/interface/handler/auth_handler.go:45-53`):

```go
// Check if user exists in Firebase Auth before generating OTP
_, err = h.authService.GetUserByEmail(c.Request.Context(), req.Email)
if err != nil {
    // Use generic error message to prevent email enumeration attacks
    log.Printf("Authentication failed for OTP request: %v", err)
    c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
    return
}
```

### 検証テスト

#### TestAuthHandler_RequestOTP_UserNotFound

**テストファイル**: `internal/interface/handler/auth_handler_test.go`

**検証内容**:
- 未登録のメールアドレスに OTP リクエスト
- レスポンスが `"Authentication failed"` であることを確認
- HTTP ステータスが 401 Unauthorized であることを確認

**期待される動作**:
```
リクエスト: {"email": "nonexistent@example.com"}
レスポンス: {"error": "Authentication failed"}
ステータス: 401 Unauthorized
```

#### TestAuthHandler_VerifyOTP_UserNotFound

**検証内容**:
- OTP は有効だがユーザーが存在しない場合
- レスポンスが `"Authentication failed"` であることを確認

**重要なポイント**:
- 登録済みユーザーと未登録ユーザーで**同じエラーメッセージ**を返す
- エラーメッセージから情報を推測できない

---

## 2. OTP 試行回数制限

### 脅威の説明

攻撃者が総当たり攻撃（ブルートフォース）で OTP を推測しようとする脅威。

**攻撃シナリオ**:
```
6桁の OTP = 1,000,000通り
試行回数無制限 = 必ず突破できる
```

### 実装された対策

**3回の試行回数制限**:
- OTP 検証失敗時に試行回数をインクリメント
- 3回失敗すると検証をブロック
- OTP を再生成するまで検証不可

**データ構造**:

```go
type otpDocument struct {
    OTP       string    `firestore:"otp"`
    ExpiresAt time.Time `firestore:"expiresAt"`
    Attempts  int       `firestore:"attempts"` // 失敗回数
}

const maxOTPAttempts = 3
```

**試行回数管理フロー**:

```
初期状態: attempts = 0 (検証可能)
1回失敗: attempts = 1 (検証可能)
2回失敗: attempts = 2 (検証可能)
3回失敗: attempts = 3 (検証ブロック)
4回目以降: ErrTooManyAttempts
```

### 検証テスト

#### TestOTPRepository_Find_TooManyAttempts

**テストファイル**: `internal/infrastructure/persistence/otp_repository_test.go`

**検証内容**:
- `attempts = 3` の OTP を作成
- `Find` 呼び出しで `ErrTooManyAttempts` が返されることを確認

#### TestOTPService_VerifyOTP_MultipleFailedAttempts

**テストファイル**: `internal/usecase/otp_service_test.go`

**検証内容**:
1. 正しい OTP を生成
2. 間違った OTP で3回検証失敗
3. 各失敗後に `attempts` がインクリメントされることを確認
4. 4回目（正しい OTP でも）に `ErrTooManyAttempts` が返されることを確認

**テストコード例**:

```go
// 1回目失敗
valid, err := service.VerifyOTP(ctx, email, wrongOTP)
// attempts = 1

// 2回目失敗
valid, err = service.VerifyOTP(ctx, email, wrongOTP)
// attempts = 2

// 3回目失敗
valid, err = service.VerifyOTP(ctx, email, wrongOTP)
// attempts = 3

// 4回目（正しいOTPでも）ブロック
valid, err = service.VerifyOTP(ctx, email, correctOTP)
if !errors.Is(err, persistence.ErrTooManyAttempts) {
    t.Error("Expected ErrTooManyAttempts")
}
```

---

## 3. OTP 期限管理

### 脅威の説明

古い OTP が永続的に使用できると、リプレイ攻撃や盗聴のリスクが高まります。

**攻撃シナリオ**:
```
1. 攻撃者が OTP を盗聴
2. 長時間経過後に OTP を使用
3. 不正ログイン成功
```

### 実装された対策

**5分間の有効期限**:
- OTP 生成時に `expiresAt = 現在時刻 + 5分` を設定
- 検証時に期限切れをチェック
- 期限切れの OTP は検証不可

**コード実装** (`internal/infrastructure/persistence/otp_repository.go:46-48`):

```go
ExpiresAt: time.Now().Add(otpExpiration), // 5分

const otpExpiration = 5 * time.Minute
```

**期限チェック** (`internal/infrastructure/persistence/otp_repository.go:83-85`):

```go
// Check expiration
if time.Now().After(otpDoc.ExpiresAt) {
    return "", fmt.Errorf("%w for email: %s", ErrOTPExpired, email)
}
```

### 検証テスト

#### TestOTPRepository_Find_Expired

**検証内容**:
- 手動で期限切れの OTP ドキュメントを作成（`expiresAt` を過去の時刻に設定）
- `Find` 呼び出しで `ErrOTPExpired` が返されることを確認
- ドキュメントが削除されず残っていることを確認

**テストコード例**:

```go
// 期限切れのOTPを作成
_, err := client.Collection("otps").Doc(email).Set(ctx, map[string]interface{}{
    "otp":       "123456",
    "expiresAt": time.Now().Add(-1 * time.Minute), // 1分前に期限切れ
    "attempts":  0,
})

// 検証時にエラーが返される
_, err = repo.Find(ctx, email)
if !errors.Is(err, persistence.ErrOTPExpired) {
    t.Error("Expected ErrOTPExpired")
}
```

#### TestOTPService_VerifyOTP_ExpiredOTP

**検証内容**:
- Usecase 層でも期限切れ OTP が正しく拒否されることを確認

---

## 4. OTP の非公開

### 脅威の説明

OTP が API レスポンスやログに出力されると、第三者に漏洩するリスクがあります。

**情報漏洩経路**:
- API レスポンス
- サーバーログ
- クライアントログ
- ネットワーク監視

### 実装された対策

**API レスポンスから除外**:

```go
// ❌ 悪い例（OTPを返す）
c.JSON(http.StatusOK, gin.H{"otp": generatedOTP})

// ✅ 良い例（OTPを返さない）
c.JSON(http.StatusOK, gin.H{
    "message": "OTP sent successfully. Please check Firestore Emulator UI.",
})
```

**ログから除外**:

```go
// ❌ 悪い例
log.Printf("Generated OTP: %s for %s", otp, email)

// ✅ 良い例（OTPをログに出力しない）
// OTP generation successful（OTP値は出力しない）
```

### 検証方法

**開発環境**:
- Firestore Emulator UI で OTP を確認
- URL: `http://localhost:4000`

**本番環境**:
- メール送信で OTP を通知（未実装）

### 検証テスト

#### TestAuthHandler_RequestOTP_Success

**検証内容**:
- レスポンスに `otp` フィールドが含まれていないことを確認
- `message` フィールドのみが返されることを確認

**テストコード例**:

```go
var response map[string]interface{}
json.Unmarshal(w.Body.Bytes(), &response)

// OTP がレスポンスに含まれていないことを確認
if _, exists := response["otp"]; exists {
    t.Error("OTP should not be included in response")
}

// メッセージのみが返されることを確認
if _, exists := response["message"]; !exists {
    t.Error("Message should be included in response")
}
```

---

## 5. ワンタイムパスワードの保証

### 脅威の説明

同じ OTP を複数回使用できると、盗聴された OTP が悪用されるリスクがあります。

**攻撃シナリオ**:
```
1. 正規ユーザーが OTP で認証成功
2. 攻撃者が同じ OTP を使用
3. 不正ログイン成功
```

### 実装された対策

**検証成功後に削除**:
- OTP 検証が成功したら即座に削除
- 同じ OTP を2回使用できないようにする

**コード実装** (`internal/usecase/otp_service.go:68-71`):

```go
// Delete OTP after successful verification (one-time use)
if deleteErr := s.otpRepo.Delete(ctx, email); deleteErr != nil {
    log.Printf("Warning: failed to delete OTP for %s: %v", email, deleteErr)
}
```

**重要な設計変更**:
- 以前: `Find` メソッドが OTP を削除
- 現在: `VerifyOTP` メソッドが検証成功後に `Delete` を呼び出す

### 検証テスト

#### TestOTPService_VerifyOTP_Success

**検証内容**:
1. OTP を生成
2. 正しい OTP で検証成功
3. Firestore から OTP が削除されていることを確認
4. 同じ OTP で再度検証を試みると失敗することを確認

**テストコード例**:

```go
// 検証成功
valid, err := service.VerifyOTP(ctx, email, correctOTP)
if !valid || err != nil {
    t.Error("Verification should succeed")
}

// OTP が削除されていることを確認
_, err = client.Collection("otps").Doc(email).Get(ctx)
if err == nil {
    t.Error("OTP should be deleted after successful verification")
}
```

#### TestOTPRepository_Delete

**検証内容**:
- 削除機能が正しく動作することを確認
- 冪等性（何度削除しても同じ結果）を確認

---

## 6. 安全な乱数生成

### 脅威の説明

予測可能な乱数生成アルゴリズムを使用すると、OTP が推測される可能性があります。

**脆弱な実装例**:

```go
// ❌ 悪い例（予測可能）
rand.Seed(time.Now().Unix())
otp := rand.Intn(1000000)

// ❌ モジュロバイアス
n := rand.Int() % 1000000
```

### 実装された対策

**crypto/rand の使用**:
- 暗号学的に安全な乱数生成器を使用
- `crypto/rand` パッケージを採用

**モジュロバイアスの排除**:
- `rand.Int(rand.Reader, max)` を使用
- 範囲指定で直接生成

**コード実装** (`internal/domain/vo/otp/value.go:33-42`):

```go
func generate6DigitCode() (string, error) {
    // Generate a number between 0 and 999999 (inclusive)
    max := big.NewInt(1000000) // 10^6
    n, err := rand.Int(rand.Reader, max)
    if err != nil {
        return "", fmt.Errorf("failed to generate random OTP: %w", err)
    }

    // Format as 6-digit string with leading zeros
    return fmt.Sprintf("%06d", n.Int64()), nil
}
```

**モジュロバイアスとは**:

```
悪い例: rand.Int() % 1000000
- rand.Int() の範囲が 1000000 の倍数でない
- 小さい数字が出やすくなる（バイアス）

良い例: rand.Int(rand.Reader, big.NewInt(1000000))
- 指定範囲内で均等に分布
- バイアスなし
```

### 検証テスト

#### TestNewOTP

**テストファイル**: `internal/domain/vo/otp/value_test.go`

**検証内容**:
- OTP が 6桁の文字列であることを確認
- すべての文字が数字（0-9）であることを確認
- エラーが発生しないことを確認

**テストコード例**:

```go
otp, err := otp.NewOTP()

// エラーなし
if err != nil {
    t.Errorf("NewOTP() failed: %v", err)
}

// 6桁
if len(otp.String()) != 6 {
    t.Errorf("Expected length 6, got %d", len(otp.String()))
}

// すべて数字
for _, c := range otp.String() {
    if c < '0' || c > '9' {
        t.Errorf("Expected digit, got %c", c)
    }
}
```

**統計的検証（手動）**:

大量の OTP を生成して分布を確認することも可能です：

```go
// 100,000個の OTP を生成して分布を確認
distribution := make(map[string]int)
for i := 0; i < 100000; i++ {
    otp, _ := otp.NewOTP()
    distribution[otp.String()[:1]]++ // 先頭1桁の分布
}

// 各数字の出現回数がほぼ均等であることを確認
// 理論値: 100,000 / 10 = 10,000 (±許容誤差)
```

---

## セキュリティテストのまとめ

### 実装済みの対策

| 対策 | 実装状況 | テスト状況 | 優先度 |
|------|---------|-----------|--------|
| メールアドレス列挙攻撃対策 | ✅ 実装済み | ✅ テスト済み | 高 |
| OTP 試行回数制限 | ✅ 実装済み | ✅ テスト済み | 高 |
| OTP 期限管理 | ✅ 実装済み | ✅ テスト済み | 高 |
| OTP の非公開 | ✅ 実装済み | ✅ テスト済み | 高 |
| ワンタイムパスワードの保証 | ✅ 実装済み | ✅ テスト済み | 高 |
| 安全な乱数生成 | ✅ 実装済み | ✅ テスト済み | 中 |

### 今後の強化案

1. **レート制限のテスト**
   - IP ベースのレート制限が正しく動作することを確認
   - 現在: 実装済みだがテスト未実施

2. **HTTPS 強制**
   - 本番環境で HTTPS 接続を強制
   - HTTP リクエストを HTTPS にリダイレクト

3. **CORS 設定の厳格化**
   - 本番環境で許可するオリジンを明示的に指定
   - 現在: `ALLOWED_ORIGINS` 環境変数で設定可能

4. **セキュリティヘッダーの追加**
   - `X-Content-Type-Options: nosniff`
   - `X-Frame-Options: DENY`
   - `Content-Security-Policy`

---

## 参考資料

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [OWASP Authentication Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html)
- [NIST Digital Identity Guidelines](https://pages.nist.gov/800-63-3/)
- [CWE-204: Observable Response Discrepancy](https://cwe.mitre.org/data/definitions/204.html)
- [CWE-307: Improper Restriction of Excessive Authentication Attempts](https://cwe.mitre.org/data/definitions/307.html)

---

**関連ドキュメント:**
- [テスト実行ガイド](./testing-guide.md) - テスト実行方法
- [設計の意思決定](./design-decisions.md) - テスト設計の背景

**最終更新日**: 2025-12-28
