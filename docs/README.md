# GBF Discord Bot ドキュメント

このディレクトリには、GBF (グランブルーファンタジー) Discord Bot の包括的なドキュメントが含まれています。

## 📚 ドキュメント構成

### 👨‍💻 エンジニア・開発者向け (`technical/`)

技術的な詳細、アーキテクチャ、開発環境について説明しています。

- **[システム設計書](technical/architecture.md)** - システム全体のアーキテクチャと構成要素
- **[API仕様書](technical/api-reference.md)** - Discord コマンドと内部APIの詳細仕様
- **[開発者ガイド](technical/development-guide.md)** - 開発環境のセットアップと開発手順
- **[データベース設計](technical/database-design.md)** - データベーススキーマと関係性
- **[移植ガイド](technical/migration-guide.md)** - Python版からGo版への移植内容と差分

### 👥 利用者・管理者向け (`user/`)

Bot の使用方法、設定、管理について説明しています。

- **[利用者マニュアル](user/user-manual.md)** - コマンドの使い方と機能一覧
- **[セットアップガイド](user/setup-guide.md)** - Discord サーバーへの Bot 導入方法
- **[管理者ガイド](user/admin-guide.md)** - サーバー管理者向けの設定と運用方法
- **[FAQ・トラブルシューティング](user/faq.md)** - よくある質問と問題解決

### 🎯 機能詳細

- **[バトル募集システム](technical/battle-recruitment.md)** - マルチバトル募集機能の詳細
- **[スケジュール管理](technical/schedule-management.md)** - イベント通知とスケジュール機能
- **[外部連携](technical/external-integrations.md)** - Google Spreadsheet 連携など

## 🚀 クイックスタート

### 利用者の方

1. [利用者マニュアル](user/user-manual.md) でコマンドの使い方を確認
2. [FAQ](user/faq.md) で疑問を解決

### サーバー管理者の方

1. [セットアップガイド](user/setup-guide.md) で Bot を導入
2. [管理者ガイド](user/admin-guide.md) で詳細設定を実施

### 開発者の方

1. [開発者ガイド](technical/development-guide.md) で開発環境を構築
2. [システム設計書](technical/architecture.md) で全体構成を把握
3. [API仕様書](technical/api-reference.md) で実装詳細を確認

## 📖 Python版からの変更点

このGo版 Bot は、元のPython版 Bot の機能を維持しながら、以下の改善を行っています：

- **パフォーマンス向上**: Go言語による高速処理
- **メモリ効率**: 軽量なリソース使用量
- **保守性向上**: 構造化されたコードアーキテクチャ
- **テスト充実**: 包括的な単体テスト・統合テスト

詳細は [移植ガイド](technical/migration-guide.md) をご覧ください。

## 🔗 関連リンク

- [プロジェクト GitHub リポジトリ](https://github.com/varubogu/gbf_discord_bot_go)
- [Python版リポジトリ](https://github.com/varubogu/gbf_bot)
- [Discord Developer Portal](https://discord.com/developers/applications)

## 📝 ドキュメント更新履歴

- **2025-08-18**: 初版作成（Go版移植完了に伴う包括的ドキュメント整備）

---

**Note**: このドキュメントは継続的に更新されます。最新情報は GitHub リポジトリをご確認ください。