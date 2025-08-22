# GBF Discord Bot Go

Go製のグラブル支援Discord Bot - discord.pyからの移植版

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## 概要

このプロジェクトは、Python（discord.py）で作成されたグラブル支援Botを Go言語 + discordgo に移植したものです。
高性能で軽量なDiscord Botとして、グランブルファンタジーに関連する機能を提供します。

## 特徴

- 🚀 高性能・軽量（Goベース）
- 📝 構造化ログ（JSON形式）
- 🔧 環境変数による設定管理
- 🎯 Prefix コマンドと Slash コマンド対応
- 🧪 包括的なテストカバレッジ
- 🔒 最小権限設計

## 前提条件

- Go 1.21 以上（推奨: 1.22+）
- Discord Developer Portal でのBot作成
- （オプション）GBF API アクセス

## インストール

### 1. リポジトリのクローン

```bash
git clone https://github.com/varubogu/gbf_discord_bot_go.git
cd gbf_discord_bot_go
```

### 2. 依存関係のインストール

```bash
go mod download
```

### 3. 環境設定

```bash
# .env ファイルを作成
cp .env.example .env
# エディタで .env を編集し、適切な値を設定
```

## 設定

### 必須環境変数

| 変数名 | 説明 | 例 |
|--------|------|-----|
| `DISCORD_TOKEN` | Discord Bot Token | `your_discord_bot_token_here` |

### オプション環境変数

| 変数名 | デフォルト | 説明 |
|--------|------------|------|
| `LOG_LEVEL` | `info` | ログレベル（debug/info/warn/error） |
| `TEST_GUILD_ID` | - | テスト用ギルドID（開発時） |
| `TEST_CHANNEL_ID` | - | テスト用チャンネルID（開発時） |

### Discord Developer Portal 設定

1. [Discord Developer Portal](https://discord.com/developers/applications) でアプリケーションを作成
2. Bot タブでボットを作成し、トークンを取得
3. **Privileged Gateway Intents** で以下を有効化：
   - Message Content Intent（Prefix コマンド使用時）
4. OAuth2 > URL Generator で適切な権限を設定してBotを招待

## 📚 ドキュメント

このプロジェクトには包括的なドキュメントが用意されています：

### 👥 利用者向けドキュメント
- **[利用者マニュアル](docs/user/user-manual.md)** - コマンドの使い方と機能一覧
- **[セットアップガイド](docs/user/setup-guide.md)** - Discord サーバーへの Bot 導入方法
- **[管理者ガイド](docs/user/admin-guide.md)** - サーバー管理者向けの設定と運用方法
- **[FAQ・トラブルシューティング](docs/user/faq.md)** - よくある質問と問題解決

### 👨‍💻 開発者向けドキュメント
- **[システム設計書](docs/technical/architecture.md)** - システム全体のアーキテクチャと構成要素
- **[API仕様書](docs/technical/api-reference.md)** - Discord コマンドと内部APIの詳細仕様

### 🚀 クイックスタート
- **利用者の方**: [利用者マニュアル](docs/user/user-manual.md) から始めてください
- **サーバー管理者の方**: [セットアップガイド](docs/user/setup-guide.md) で Bot を導入
- **開発者の方**: [システム設計書](docs/technical/architecture.md) で全体構成を把握

## 使用方法

### 起動

```bash
# ビルドして実行
go build -o gbf-bot ./cmd/bot
./gbf-bot

# または直接実行
go run ./cmd/bot
```

### 主要機能

#### 基本コマンド
- `!ping` / `/ping` - Bot の応答確認
- `!help` / `/help` - コマンド一覧と使用方法
- `!status` / `/status` - Bot の状態表示

#### バトル関連機能
- `!battles` / `/battles` - 利用可能なバトル一覧
- `!battle <id>` / `/battle <id>` - バトル詳細情報

#### 管理機能（要権限）
- `!reload` / `/reload` - Bot設定の再読み込み

詳細なコマンド仕様は [利用者マニュアル](docs/user/user-manual.md) をご覧ください。

## 開発

### ディレクトリ構成

```
├── cmd/bot/           # エントリポイント
├── internal/
│   ├── config/        # 設定管理
│   ├── discord/       # Discord セッション・ハンドラ
│   ├── commands/      # コマンド実装
│   ├── gbf/           # GBF ドメインロジック
│   └── log/           # ログ設定
├── testdata/          # テスト用データ
├── .env.example       # 環境変数テンプレート
├── .golangci.yml      # リンター設定
└── README.md
```

### テストの実行

```bash
# 全てのテスト実行
go test ./...

# カバレッジ付きテスト
go test -cover ./...

# 統合テスト実行（環境変数設定必要）
go test -tags=integration ./...
```

### コードの品質チェック

```bash
# フォーマット
go fmt ./...

# Vet チェック
go vet ./...

# Lint（golangci-lint インストール必要）
golangci-lint run
```

### ビルド

```bash
# 開発用ビルド
go build ./cmd/bot

# 本番用ビルド（最適化）
CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o gbf-bot ./cmd/bot

# クロスプラットフォームビルド
GOOS=windows GOARCH=amd64 go build -o gbf-bot.exe ./cmd/bot
GOOS=darwin GOARCH=amd64 go build -o gbf-bot-darwin ./cmd/bot
```

## 移植ステータス

### ✅ 完了

- [x] 基本プロジェクト構造
- [x] 設定管理（環境変数）
- [x] 構造化ログ
- [x] Discord セッション初期化
- [x] 基本的なイベントハンドラ
- [x] Ping コマンド（Prefix + Slash）
- [x] テスト基盤
- [x] リンター設定

### 🚧 進行中

- [ ] 追加コマンド実装
- [ ] GBF API 連携
- [ ] UI コンポーネント（Button, Select）
- [ ] 定期タスク機能

### 📋 未実装

- [ ] 完全な機能パリティ（../gbf_bot から）
- [ ] データベース連携
- [ ] キャッシュ機能
- [ ] 高度なエラーハンドリング
- [ ] CI/CD パイプライン

## 貢献

1. フォークを作成
2. フィーチャーブランチを作成（`git checkout -b feature/amazing-feature`）
3. 変更をコミット（`git commit -m 'Add amazing feature'`）
4. ブランチをプッシュ（`git push origin feature/amazing-feature`）
5. Pull Request を作成

## トラブルシューティング

### よくある問題

**Q: Bot が応答しない**
- Discord Token が正しく設定されているか確認
- Bot に適切な権限が付与されているか確認
- Message Content Intent が有効になっているか確認（Prefix コマンド使用時）

**Q: ログレベルを変更したい**
- `LOG_LEVEL` 環境変数を設定（debug/info/warn/error）

**Q: Slash コマンドが表示されない**
- `TEST_GUILD_ID` を設定してテストギルドで確認
- Discord への伝搬には最大1時間かかる場合があります

## ライセンス

このプロジェクトは MIT ライセンスの下で公開されています。詳細は [LICENSE](LICENSE) ファイルを参照してください。

## 関連リンク

- [Discord.js Guide](https://discordjs.guide/)
- [discordgo Documentation](https://pkg.go.dev/github.com/bwmarrin/discordgo)
- [Go Documentation](https://golang.org/doc/)
