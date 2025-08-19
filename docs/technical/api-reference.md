# API仕様書

GBF Discord Bot (Go版) のコマンドAPI仕様と内部関数について説明します。

## 📋 目次

- [コマンドAPI概要](#コマンドapi概要)
- [基本コマンド](#基本コマンド)
- [バトル関連コマンド](#バトル関連コマンド)
- [管理コマンド](#管理コマンド)
- [内部API](#内部api)
- [エラーハンドリング](#エラーハンドリング)
- [レスポンス形式](#レスポンス形式)

## 🎯 コマンドAPI概要

### サポート形式
- **Prefix Commands**: `!command [args...]`
- **Slash Commands**: `/command [options...]`

### 共通仕様
- **認証**: Discord OAuth2経由
- **権限**: ロールベース制御
- **ログ**: 構造化ログによる実行記録
- **エラー**: ユーザーフレンドリーなエラーメッセージ

### レスポンス時間
- **即座応答**: ping, help等の軽量コマンド
- **遅延応答**: 複雑な処理が必要なコマンド（3秒制限対応）

## 🎮 基本コマンド

### ping

Botの応答確認とレイテンシ測定を行います。

#### Prefix Command
```
!ping
```

#### Slash Command
```
/ping
```

#### パラメータ
なし

#### レスポンス例
```
Pong!
```

#### 実装詳細
- **ファイル**: `internal/commands/ping.go`
- **権限**: なし（全ユーザー利用可）
- **処理時間**: < 100ms
- **ログ**: 実行ログとコンテキスト情報記録

---

### help

利用可能なコマンドの一覧と使用方法を表示します。

#### Prefix Command
```
!help [command_name]
```

#### Slash Command
```
/help [command: command_name]
```

#### パラメータ
| 名前 | 型 | 必須 | 説明 |
|------|----|----|------|
| command | string | No | 詳細を表示する特定のコマンド名 |

#### レスポンス例

**一般ヘルプ:**
```markdown
# GBF Discord Bot Help

## General
- `!ping` / `/ping` - Pings the bot and returns response time
- `!help` / `/help` - Shows this help message

## GBF
- `!battles` / `/battles` - Shows list of available battles
- `!battle <id>` / `/battle <id>` - Shows detailed information about a specific battle
```

**特定コマンドヘルプ:**
```markdown
# Help: ping

Pings the bot and returns response time

**Usage:** !ping or /ping
**Category:** General
**Available as:** Prefix (!), Slash (/)
```

#### 実装詳細
- **ファイル**: `internal/commands/help.go`
- **権限**: なし（全ユーザー利用可）
- **機能**: 動的コマンド一覧生成、カテゴリ別表示

---

## ⚔️ バトル関連コマンド

### battles

利用可能なバトル一覧を表示します。

#### Prefix Command
```
!battles [type]
```

#### Slash Command
```
/battles [type: battle_type]
```

#### パラメータ
| 名前 | 型 | 必須 | 説明 |
|------|----|----|------|
| type | string | No | フィルターするバトルタイプ |

#### 利用可能なバトルタイプ
- `hl` - High Level raids
- `faa_hl` - Lucilius Hard
- `baha_hl` - Proto Bahamut Hard
- `ubaha_hl` - Ultimate Bahamut Hard
- `akasha_hl` - Akasha Hard
- `luci_hl` - Lucifer Hard
- `gw` - Guild War
- `gw_nm` - Guild War Nightmare
- `event` - Event raids
- `event_hl` - Event High Level
- `train` - Training rooms
- `custom` - Custom battles

#### レスポンス例
```markdown
# GBF Battle List - Active Battles

## HL Battles
🟢 `faa_hl` - Lucilius (Hard) (Lv.200)
🟢 `baha_hl` - Proto Bahamut (Hard) (Lv.150)
🟢 `ubaha_hl` - Ultimate Bahamut (Hard) (Lv.200)

Use !battle <id> or /battle <id> for detailed information
```

#### 実装詳細
- **ファイル**: `internal/commands/battle.go`
- **ドメインロジック**: `internal/gbf/battle.go`
- **権限**: なし（全ユーザー利用可）

---

### battle

特定バトルの詳細情報を表示します。

#### Prefix Command
```
!battle <battle_id>
```

#### Slash Command
```
/battle id:<battle_id>
```

#### パラメータ
| 名前 | 型 | 必須 | 説明 |
|------|----|----|------|
| id | string | Yes | 詳細を表示するバトルID |

#### レスポンス例
```markdown
# Battle Info: Lucilius (Hard)

Dark Rapture Hard mode raid

**Battle ID:** faa_hl
**Level:** 200
**Type:** faa_hl
**Min Rank:** 150
**Max Players:** 6
**Status:** 🟢 Active

Created: 2025-08-18 12:00:00
```

#### エラーレスポンス例
```markdown
# Battle Not Found

No battle found with ID: `invalid_id`
```

#### 実装詳細
- **ファイル**: `internal/commands/battle.go`
- **ドメインロジック**: `internal/gbf/battle.go`
- **権限**: なし（全ユーザー利用可）

---

## 🛠️ 管理コマンド

### reload

Bot設定とコンポーネントを再読み込みします。

#### Prefix Command
```
!reload
```

#### Slash Command
```
/reload
```

#### パラメータ
なし

#### レスポンス例
**成功時:**
```
✅ Bot components reloaded successfully!
```

**権限不足時:**
```
❌ You don't have permission to use this command. Required role: `gbf_bot_control`
```

#### 実装詳細
- **ファイル**: `internal/commands/admin.go`
- **権限**: `gbf_bot_control` ロールまたは管理者権限
- **処理内容**: 
  - 設定再読み込み
  - 外部接続のリフレッシュ
  - コマンド登録の更新
  - キャッシュのクリア

---

### status

Botの現在状態を表示します。

#### Prefix Command
```
!status
```

#### Slash Command
```
/status
```

#### パラメータ
なし

#### レスポンス例
```markdown
# Bot Status

**Status:** ✅ Online
**Version:** v0.1.0
**Language:** Go
**Commands:** ping, help, reload, status
```

#### 実装詳細
- **ファイル**: `internal/commands/admin.go`
- **権限**: なし（全ユーザー利用可）
- **情報源**: 
  - システム状態
  - バージョン情報
  - 利用可能コマンド一覧

---

## 🔧 内部API

### BattleManager

バトル情報の管理を行うドメインロジック。

#### 主要メソッド

```go
type BattleManager struct {
    battles map[string]*BattleInfo
}

func NewBattleManager() *BattleManager
func (bm *BattleManager) GetBattle(id string) (*BattleInfo, error)
func (bm *BattleManager) GetActiveBattles() []*BattleInfo
func (bm *BattleManager) GetBattlesByType(battleType BattleType) []*BattleInfo
```

#### BattleInfo構造体
```go
type BattleInfo struct {
    ID          string
    Name        string
    Type        BattleType
    Level       int
    MinRank     int
    MaxPlayers  int
    Description string
    IsActive    bool
    CreatedAt   time.Time
}
```

---

### RecruitmentManager

募集管理のドメインロジック（実装予定）。

#### 主要メソッド
```go
type RecruitmentManager struct {
    recruitments map[string]*Recruitment
    battleManager *BattleManager
}

func NewRecruitmentManager(battleManager *BattleManager) *RecruitmentManager
func (rm *RecruitmentManager) CreateRecruitment(req *Recruitment) error
func (rm *RecruitmentManager) GetRecruitment(id string) (*Recruitment, error)
func (rm *RecruitmentManager) AddParticipant(recruitmentID, userID, username string) error
```

---

### AttackCalculator

GBF関連の計算処理。

#### 主要メソッド
```go
type AttackCalculator struct{}

func NewAttackCalculator() *AttackCalculator
func (c *AttackCalculator) CalculateBaseDamage(attack int, defense int) int
func (c *AttackCalculator) CalculateCriticalDamage(baseDamage int, critMultiplier float64) int
func (c *AttackCalculator) CalculateElementalDamage(baseDamage int, elementalModifier float64) int
```

---

## ⚠️ エラーハンドリング

### エラー分類

#### 1. ユーザーエラー (400番台相当)
- **権限不足**: 必要なロール/権限がない
- **不正な入力**: 無効なパラメータやID
- **リソース不足**: 募集枠が満員等

#### 2. システムエラー (500番台相当)
- **内部処理エラー**: 予期しない処理失敗
- **外部サービスエラー**: Discord API、データベース接続等
- **タイムアウト**: 処理時間超過

### エラーレスポンス形式

#### ユーザーエラー
```
❌ [エラーメッセージ]
💡 [解決方法の提案]
```

例:
```
❌ You don't have permission to use this command. Required role: `gbf_bot_control`
💡 Ask a server administrator to assign you the required role.
```

#### システムエラー
```
🚨 An error occurred while processing your command.
🔧 Please try again later. If the problem persists, contact an administrator.
```

### ログ記録

#### 成功ログ
```json
{
  "level": "info",
  "msg": "Command executed successfully",
  "guild_id": "123456789",
  "channel_id": "987654321",
  "user_id": "456789123",
  "command": "ping",
  "execution_time_ms": 45
}
```

#### エラーログ
```json
{
  "level": "error",
  "msg": "Command execution failed",
  "guild_id": "123456789",
  "channel_id": "987654321",
  "user_id": "456789123",
  "command": "battle",
  "error": "battle not found: invalid_id"
}
```

---

## 📊 レスポンス形式

### Embed形式

バトル情報や複雑な応答には Discord Embed を使用。

#### 基本構造
```json
{
  "title": "タイトル",
  "description": "説明文",
  "color": 0x3498db,
  "fields": [
    {
      "name": "フィールド名",
      "value": "値",
      "inline": true
    }
  ],
  "footer": {
    "text": "フッター情報"
  }
}
```

### プレーンテキスト形式

シンプルなコマンドではプレーンテキストを使用。

#### 例
```
Pong!
✅ Bot components reloaded successfully!
❌ Battle not found: invalid_id
```

---

## 🔄 将来の拡張予定

### 募集コマンド
```
/recruit quest:<quest_name> [battle_type:<type>] [time:<datetime>]
/recruitment list
/recruitment join <recruitment_id>
/recruitment leave <recruitment_id>
```

### スケジュールコマンド
```
/schedule list
/schedule add <event_name> <datetime>
/schedule remove <event_id>
```

### 統計コマンド
```
/stats battles
/stats recruitment
/stats user [user_id]
```

---

**関連ドキュメント:**
- [システム設計書](architecture.md) - システム全体のアーキテクチャ
- [開発者ガイド](development-guide.md) - 開発環境とワークフロー
- [利用者マニュアル](../user/user-manual.md) - エンドユーザー向け使用方法

**更新履歴:**
- 2025-08-18: 初版作成（基本コマンド・バトル関連・管理コマンド）

**Note:** この仕様書は実装状況に応じて継続的に更新されます。新機能追加時は対応する API 仕様も更新してください。