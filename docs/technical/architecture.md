# システム設計書

GBF Discord Bot (Go版) のシステムアーキテクチャと技術仕様について説明します。

## 📋 目次

- [システム概要](#システム概要)
- [アーキテクチャ](#アーキテクチャ)
- [技術スタック](#技術スタック)
- [ディレクトリ構成](#ディレクトリ構成)
- [コンポーネント設計](#コンポーネント設計)
- [データフロー](#データフロー)
- [セキュリティ設計](#セキュリティ設計)
- [パフォーマンス設計](#パフォーマンス設計)

## 🎯 システム概要

### プロジェクト目標
- Python版GBF Discord Botの機能をGoで再実装
- パフォーマンス向上とリソース効率の改善
- 保守性とテスト容易性の向上
- スケーラビリティの確保

### 主要機能
1. **Discord Bot 基盤**
   - Slash Commands と Prefix Commands 両対応
   - イベント処理とリアクション管理
   - 権限管理とセキュリティ

2. **バトル募集システム**
   - マルチバトル募集投稿
   - 参加者管理とリアクション処理
   - 自動期限切れ処理

3. **スケジュール・通知システム**
   - 古戦場・イベント通知
   - カスタムスケジュール管理
   - 定期タスク実行

4. **外部連携**
   - Google Spreadsheet連携（計画中）
   - データベース連携（計画中）

## 🏗️ アーキテクチャ

### 全体アーキテクチャ図

```
┌─────────────────────────────────────────────────────────┐
│                    Discord API                          │
└─────────────────────┬───────────────────────────────────┘
                      │ WebSocket/REST
┌─────────────────────▼───────────────────────────────────┐
│                  Discord Gateway                        │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│                Application Layer                        │
├─────────────────────────────────────────────────────────┤
│  cmd/bot/main.go                                        │
│  ├── Configuration Loading                              │
│  ├── Logger Initialization                              │
│  ├── Bot Instance Creation                              │
│  └── Graceful Shutdown                                  │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│                  Bot Core Layer                         │
├─────────────────────────────────────────────────────────┤
│  internal/discord/bot.go                                │
│  ├── Session Management                                 │
│  ├── Event Handler Registration                         │
│  ├── Command Routing                                    │
│  └── Lifecycle Management                               │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│                Command Layer                            │
├─────────────────────────────────────────────────────────┤
│  internal/commands/                                     │
│  ├── ping.go     - Health Check Commands               │
│  ├── help.go     - Help & Documentation                │
│  ├── admin.go    - Administrative Commands             │
│  ├── battle.go   - Battle Information Commands         │
│  └── (future)    - Recruitment Commands                │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│                 Domain Layer                            │
├─────────────────────────────────────────────────────────┤
│  internal/gbf/                                          │
│  ├── calc.go        - Battle Calculations              │
│  ├── battle.go      - Battle Management                │
│  ├── recruitment.go - Recruitment Management           │
│  └── (future)       - Schedule Management              │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│               Infrastructure Layer                      │
├─────────────────────────────────────────────────────────┤
│  internal/config/  - Configuration Management          │
│  internal/log/     - Structured Logging                │
│  └── (future)      - Database, External APIs           │
└─────────────────────────────────────────────────────────┘
```

### レイヤー設計原則

#### 1. Application Layer
- **責任**: アプリケーションのエントリーポイント
- **依存**: Infrastructure Layer のみ
- **特徴**: 薄いレイヤー、設定とライフサイクル管理

#### 2. Bot Core Layer
- **責任**: Discord API との通信、イベント管理
- **依存**: Command Layer, Infrastructure Layer
- **特徴**: Discord 固有のロジック集約

#### 3. Command Layer
- **責任**: コマンド処理、ユーザーインタフェース
- **依存**: Domain Layer, Infrastructure Layer
- **特徴**: 各コマンドの独立性確保

#### 4. Domain Layer
- **責任**: ビジネスロジック、ドメイン知識
- **依存**: なし（純粋関数中心）
- **特徴**: テストしやすい、再利用可能

#### 5. Infrastructure Layer
- **責任**: 外部システムとの連携、技術的関心事
- **依存**: 外部ライブラリのみ
- **特徴**: インターフェース化、モック対応

## 💻 技術スタック

### 言語・フレームワーク
```go
// Core Language
Go 1.21+ (推奨 1.22+)

// Discord Integration
github.com/bwmarrin/discordgo v0.28.1

// Logging
log/slog (標準ライブラリ)

// Testing
標準ライブラリ + Testify (optional)

// Configuration
標準ライブラリ (os.Getenv)
```

### 外部システム（計画中）
```yaml
Database:
  - PostgreSQL 15+ (予定)
  - GORM v2 (予定)

External APIs:
  - Google Sheets API v4 (予定)
  - Discord API v10 (使用中)

Infrastructure:
  - Docker & Docker Compose
  - GitHub Actions (CI/CD)
```

### 開発ツール
```yaml
Code Quality:
  - golangci-lint
  - go fmt / go vet
  - staticcheck

Version Control:
  - Git
  - GitHub

Deployment:
  - Docker
  - Cross-platform builds
```

## 📁 ディレクトリ構成

```
gbf_discord_bot_go/
├── cmd/
│   └── bot/
│       └── main.go              # アプリケーションエントリーポイント
├── internal/
│   ├── commands/                # コマンド処理層
│   │   ├── ping.go
│   │   ├── help.go
│   │   ├── admin.go
│   │   └── battle.go
│   ├── discord/                 # Discord API 連携層
│   │   └── bot.go
│   ├── gbf/                     # ドメインロジック層
│   │   ├── calc.go
│   │   ├── calc_test.go
│   │   ├── battle.go
│   │   └── recruitment.go
│   ├── config/                  # 設定管理層
│   │   ├── config.go
│   │   └── config_test.go
│   └── log/                     # ログ管理層
│       └── log.go
├── docs/                        # ドキュメント
│   ├── README.md
│   ├── technical/
│   └── user/
├── .github/
│   └── workflows/
│       └── ci.yml               # CI/CD設定
├── go.mod                       # Go モジュール定義
├── go.sum                       # 依存関係ハッシュ
├── .env.example                 # 環境変数テンプレート
├── .golangci.yml               # 静的解析設定
└── README.md                    # プロジェクト概要
```

### ディレクトリ設計原則

#### `cmd/` - Applications
- 実行可能ファイルのエントリーポイント
- 薄いレイヤー、主に設定と起動処理
- 複数のアプリケーションに対応可能

#### `internal/` - Private Packages
- アプリケーション固有のロジック
- 外部からのimportを防ぐ
- レイヤー別にパッケージを分離

#### `docs/` - Documentation
- 利用者向け・開発者向けドキュメント
- Markdownベース
- 継続的な更新が前提

## 🔧 コンポーネント設計

### Bot Core (`internal/discord/bot.go`)

```go
type Bot struct {
    session       *discordgo.Session    // Discord セッション
    config        *config.Config        // 設定
    logger        *log.Logger          // ロガー
    pingCommand   *commands.PingCommand // コマンドハンドラー
    helpCommand   *commands.HelpCommand
    adminCommand  *commands.AdminCommand
    battleCommand *commands.BattleCommand
}
```

**設計原則:**
- 単一責任原則: Discord API との通信に集中
- 依存性注入: 各コマンドをコンストラクタで注入
- ライフサイクル管理: 起動・停止・グレースフルシャットダウン

### Command System

```go
// コマンドインターフェース（暗黙的）
type Command interface {
    HandlePrefixCommand(s *discordgo.Session, m *discordgo.MessageCreate)
    HandleSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate)
    GetSlashCommandDefinition() *discordgo.ApplicationCommand
}
```

**設計原則:**
- コマンドの独立性: 各コマンドは独立したファイル
- 両形式対応: Prefix (!cmd) と Slash (/cmd) の両方
- 統一インターフェース: 同じパターンで実装

### Domain Logic (`internal/gbf/`)

```go
// 純粋関数中心の設計
type AttackCalculator struct{}

func (c *AttackCalculator) CalculateBaseDamage(attack, defense int) int
func IsValidWeaponType(weaponType WeaponType) bool

// 状態を持つオブジェクトもインターフェース化
type BattleManager struct {
    battles map[string]*BattleInfo
}

type RecruitmentManager struct {
    recruitments map[string]*Recruitment
    battleManager *BattleManager
}
```

**設計原則:**
- 純粋関数優先: 副作用のない関数を中心に
- テストしやすさ: モックなしでテスト可能
- ドメイン知識の集約: GBF固有の知識をここに集約

### Infrastructure Layer

```go
// 設定管理
type Config struct {
    DiscordToken string
    LogLevel     string
    // 将来的に追加
    DatabaseURL  string
    GoogleKeyFile string
}

// 構造化ログ
type Logger struct {
    *slog.Logger
}

func (l *Logger) WithDiscordContext(guildID, channelID, userID string) *Logger
```

**設計原則:**
- 外部依存の抽象化: インターフェース経由でアクセス
- 設定の集約: 環境変数の一元管理
- 可観測性: 構造化ログでコンテキスト情報を付与

## 🌊 データフロー

### コマンド処理フロー

```
1. Discord Event
   ↓
2. Bot.onMessageCreate / Bot.onInteractionCreate
   ↓
3. Command Router (switch文)
   ↓
4. Specific Command Handler
   ↓
5. Domain Logic Execution
   ↓
6. Response Generation
   ↓
7. Discord API Response
```

### エラーハンドリングフロー

```go
// 統一されたエラーハンドリングパターン
func (cmd *SomeCommand) HandleCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
    logger := cmd.logger.WithDiscordContext(m.GuildID, m.ChannelID, m.Author.ID).WithCommand("command")
    
    result, err := cmd.executeLogic()
    if err != nil {
        logger.WithError(err).Error("Command execution failed")
        // ユーザーフレンドリーなメッセージを送信
        return
    }
    
    // 成功時の処理
    logger.Info("Command executed successfully")
}
```

## 🔐 セキュリティ設計

### 権限管理

```go
// 権限チェックの統一パターン
func (a *AdminCommand) hasAdminPermission(s *discordgo.Session, guildID, userID string) bool {
    // 1. 管理者権限チェック
    // 2. gbf_bot_control ロールチェック
    // 3. ログ記録
    return authorized
}
```

### 入力検証

```go
// 入力値の検証とサニタイゼーション
func validateBattleID(battleID string) error {
    if len(battleID) == 0 {
        return errors.New("battle ID cannot be empty")
    }
    // 追加の検証ロジック
    return nil
}
```

### ログ設計

```go
// セキュリティログの記録
logger.WithDiscordContext(guildID, channelID, userID).
    WithCommand("admin").
    Info("Admin command executed successfully")

logger.WithError(err).
    WithUser(userID).
    Warn("Unauthorized access attempt")
```

## ⚡ パフォーマンス設計

### メモリ管理

```go
// 効率的なデータ構造の使用
type BattleManager struct {
    battles map[string]*BattleInfo // O(1) アクセス
}

// メモリプールパターン（将来的に実装）
var embedPool = sync.Pool{
    New: func() interface{} {
        return &discordgo.MessageEmbed{}
    },
}
```

### 並行処理

```go
// Graceful Shutdown パターン
func (b *Bot) Start(ctx context.Context) error {
    // Context を使用した適切な停止処理
    select {
    case <-ctx.Done():
        b.logger.Info("Bot stopping due to context cancellation")
    case <-stop:
        b.logger.Info("Bot stopping due to interrupt signal")
    }
    return b.Close()
}
```

### レスポンス性能

```go
// Discord の3秒制限を考慮した処理
func (cmd *Command) HandleSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
    // 即座に Acknowledge
    err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
    })
    
    // 長時間処理は goroutine で実行
    go func() {
        result := cmd.processLongTask()
        s.FollowupMessageCreate(i.Interaction, &discordgo.WebhookParams{
            Content: result,
        })
    }()
}
```

## 🔄 拡張性設計

### プラグインアーキテクチャ（将来）

```go
// コマンドプラグインインターフェース
type CommandPlugin interface {
    Name() string
    Register(bot *Bot) error
    Unregister(bot *Bot) error
}

// 動的コマンド登録
func (b *Bot) RegisterPlugin(plugin CommandPlugin) error {
    return plugin.Register(b)
}
```

### 設定の外部化

```go
// 設定の段階的拡張
type Config struct {
    // 現在実装済み
    DiscordToken string
    LogLevel     string
    
    // 将来の拡張
    DatabaseURL   string
    RedisURL      string
    GoogleKeyFile string
    Features      FeatureFlags
}

type FeatureFlags struct {
    EnableRecruitment bool
    EnableSchedule    bool
    EnableGSpread     bool
}
```

---

**関連ドキュメント:**
- [開発者ガイド](development-guide.md) - 開発環境とワークフロー
- [API仕様書](api-reference.md) - コマンド仕様とAPI
- [移植ガイド](migration-guide.md) - Python版からの変更点

**更新履歴:**
- 2025-08-18: 初版作成（Go版アーキテクチャ設計完了）

**Note:** このシステムは継続的に進化します。新機能の追加に伴い、アーキテクチャも適切に更新していきます。