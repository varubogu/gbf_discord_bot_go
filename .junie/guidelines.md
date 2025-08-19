# 開発ガイドライン（gbf_discord_bot_go）

このリポジトリは、discord.py で作成された Bot（../gbf_bot）を Go へ移植する前提のプロジェクトです。現時点ではコードは未整備（README と LICENSE のみ）です。本ガイドは、移植および今後の開発に必要なプロジェクト固有の注意点・手順をまとめたものです。

---

## 1. ビルド/設定（Go + discordgo 前提）

- 推奨環境
  - Go 1.21+（推奨 1.22+）
  - Discord クライアント: github.com/bwmarrin/discordgo
  - 任意: golangci-lint、direnv/.env、Taskfile/Makefile、zerolog もしくは slog

- 初期化（未実施の場合）
  - go mod init github.com/varubogu/gbf_discord_bot_go
  - go get github.com/bwmarrin/discordgo@latest

- 環境変数（最低限）
  - DISCORD_TOKEN: Discord Bot トークン（コード側で "Bot " を付与）
  - LOG_LEVEL: info|debug など（任意）
  - GBF_API_BASE / GBF_API_KEY: 外部 GBF サービスと連携する場合に使用（任意）

- Discord Developer Portal 設定
  - 必要 Intents を有効化（Message Content Intent は prefix コマンドで必要になり得ます）
  - 可能なら Slash Command（アプリケーションコマンド）主体へ移行して Message Content 依存を低減

- 推奨ディレクトリ構成（移植時の雛形）
  - cmd/bot/main.go — エントリポイント
  - internal/discord — セッション初期化、ハンドラ、コマンド登録
  - internal/commands — 各コマンド実装（prefix/slash のルーティング源）
  - internal/gbf — グラブル関連のドメインロジック（API、計算、パース等）
  - internal/config — 環境変数ロードと検証
  - internal/log — ロギング初期化
  - testdata — テスト用フィクスチャ

- 最小起動例（要約）
  - discordgo.New("Bot "+token) でセッション作成
  - Identify.Intents に必要な Intent を設定
  - AddHandler で MessageCreate/InteractionCreate を登録
  - Open/Close で接続の開始/終了、OS シグナルでグレースフルシャットダウン

---

## 2. 移植方針（discord.py → Go）

- 機能パリティを維持しつつ、Discord 依存とドメインロジックを分離。
- discord.py の概念対応
  - Bot/commands.Bot → discordgo.Session + 独自ルーター
  - Cog → Go のパッケージ/構造体（internal/commands など）
  - @bot.command → prefix ルーターで解析（MessageCreate）
  - app_commands（Slash） → ApplicationCommand の登録と InteractionCreate で処理
  - @tasks.loop → goroutine + time.Ticker + context で制御
  - Embed / UI コンポーネント → MessageEmbed / MessageComponent（Button, Select 等）

- 段階的移行ステップ
  1) ../gbf_bot の機能棚卸し（コマンド一覧、イベント、外部依存、タスク、権限）
  2) Go 側で最小起動（ping の prefix と slash）
  3) ドメインロジックを internal/gbf に抽出移植（HTTP 呼び出し等はタイムアウトやリトライを設計）
  4) コマンドを 1 つずつ移植（共通レスポンス整備、エラー方針を統一）
  5) 背景タスク/通知系を goroutine 化（Context で停止可能に）
  6) 権限/Intents/RateLimit の最小化・遵守
  7) テストとロギングの拡充

---

## 3. テスト

- 実行コマンド
  - go test ./...（ユニットテスト）
  - go test -tags=integration ./...（統合テストは明示的に）

- 方針
  - internal/gbf などドメイン層は純粋関数/インタフェース化で単体テストしやすくする。
  - Discord 依存部は最小限のインタフェースを定義してフェイクで検証（ネットワーク不要）。
  - 統合テストはビルドタグ + 環境変数（DISCORD_TOKEN, TEST_GUILD_ID/TEST_CHANNEL_ID）で任意実行。

- 例（推奨サンプル）
  - internal/gbf/calc.go: 簡易計算関数
  - internal/gbf/calc_test.go: 上記のユニットテスト
  - 備考: 本リポジトリにはまだ作成しません。ローカルで作成して go test ./... が通ることを確認してください。標準的な内容のため Go 1.21+ でパスします。

- 追加指針
  - HTTP クライアントは interface 化し、テストでモック（レスポンス/エラー/タイムアウトの分岐を網羅）。
  - Slash/Interaction はデータを構造体に取り出して純粋ロジックを別関数に切り出すと検証が容易。

---

## 4. コードスタイル/静的解析

- go fmt / go vet は常用。
- golangci-lint 推奨（errcheck, staticcheck, gosimple, ineffassign, govet, revive 等）。
- 例の .golangci.yml（任意）
  - enable: govet, staticcheck, gosimple, ineffassign, errcheck, revive
  - run.timeout: 3m
  - revive で indent-error-flow, errorf などの基本ルール

- ロギング
  - log/slog もしくは rs/zerolog で構造化ログ。guild_id, user_id, command 等のフィールドを付与。

---

## 5. 設定・運用の注意

- 設定
  - internal/config で env を集中管理し、必須値の検証を初期化時に実施（fail fast）。
  - ローカルは .env / direnv などを利用（Secrets はコミット禁止）。

- シャットダウン/並行処理
  - SIGINT/SIGTERM を受けて context キャンセル → セッション Close → goroutine 終了。
  - goroutine 乱立を避け、背景タスクは Ticker + Context で停止可能に。

- レート制限/リトライ
  - Discord のレート制限は discordgo が一定吸収。外部 API（GBF 関連）はトークンバケットや指数バックオフを導入。

- 時刻/ロケール
  - Python 実装と挙動差が出やすい（タイムゾーン、端数処理）。time.Location や丸め方を明示し、テストで担保。

---

## 6. Slash コマンド運用

- Guild 登録は即時反映、本番の Global 登録は最大 ~1 時間の伝搬遅延。
- 開発フロー
  1) テストギルドに登録して実機確認
  2) 問題なければ Global に反映
  3) バージョン変更時はコマンド再同期の仕組み（差分更新 or 全削除→再作成）を用意

---

## 7. 典型的なハンドラ設計（要点）

- Prefix ルーター
  - 例: Router{Prefix: "!"} で MessageCreate を解析、args 分割、ハンドラへ委譲。
  - エラーはユーザー向け文言とログ文言を分離（情報漏えい防止）。

- Slash/Interaction
  - InteractionRespond/Followup を正しく使い分け（3 秒ルールに注意）。
  - 長時間処理は Deferred レスポンス + 後追い編集で返却。

- UI コンポーネント
  - CustomID をキーに状態を識別。必要なら短命トークンやメッセージ ID でスコープ管理。

---

## 8. CI/CD 推奨

- CI（GitHub Actions 想定）
  - セットアップ: actions/setup-go
  - キャッシュ: Go modules キャッシュ
  - ジョブ: golangci-lint → go test ./... → go build ./cmd/bot
  - 任意: integration タグを別ジョブ（手動/スケジュール）で実行し、DISCORD_TOKEN/TEST_CHANNEL_ID を Secrets で供給

- リリース
  - タグ v0.x 系で管理。
  - 配布するなら GOOS/GOARCH 組み合わせでクロスビルド。

---

## 9. 機能パリティ チェックリスト（移植ナビ）

- [x] 基本的な Discord Bot 設定とセッション初期化
- [x] 構造化ログシステム（slog ベース、Discord コンテキスト対応）
- [x] 環境変数管理と設定検証（internal/config）
- [x] グレースフルシャットダウンの実装
- [x] Prefix コマンド基盤（ping コマンド実装済み）
- [x] Slash Commands 基盤とテストギルド登録（ping コマンド実装済み）
- [x] イベントリスナー基本形（on_ready, on_message, on_interaction_create）
- [x] GBF ドメインロジック基盤（internal/gbf パッケージ）
- [x] 単体テスト・ベンチマークテスト（internal/gbf）
- [x] 静的解析設定（golangci-lint）
- [x] CI/CD 設定（GitHub Actions）

### 今後の拡張項目
- [ ] より多くのコマンドの移植（ヘルプ、GBF ユーティリティ等）
- [ ] Global Slash Commands への対応
- [ ] 追加イベントリスナー（on_reaction_add, on_member_* など）
- [ ] Embed/コンポーネント UI（ボタン、セレクト、モーダル）
- [ ] 定期タスクと通知（Ticker + Context）
- [ ] 外部 API 呼び出し（レート制限/再試行/キャッシュ）
- [ ] 権限システムの詳細実装
- [ ] 統合テスト環境の拡充

---

## 10. 既存テスト例に関する注意

- 本リポジトリには以下のテストファイルが実装済みです：
  - internal/gbf/calc.go: GBF関連の計算ロジック（攻撃ダメージ、武器種検証等）
  - internal/gbf/calc_test.go: 上記のユニットテスト・ベンチマークテスト
  - internal/config/config_test.go: 設定管理のテスト
- すべてのテストは `go test ./...` で実行でき、正常にパスすることを確認済みです。
- テストカバレッジは GitHub Actions の CI で自動測定され、Codecov にアップロードされます。

---

## 11. 保守メモ

- Go のバージョンや依存の更新時に本ガイドを最新化してください。
- 実際のパッケージレイアウトが固まったら、本ガイドの雛形構成を実態に合わせて更新。
- 外部サービス（GBF API 等）を導入した場合は、エンドポイント、認証方式、レート制限、タイムアウト/リトライ方針を本ガイドに追記。
