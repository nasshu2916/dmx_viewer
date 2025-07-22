## プロジェクトの概要:

このプロジェクトはバックエンドで受信した ArtNet の情報を WebSocket 経由でフロントエンドに送信し、リアルタイムで DMX データを表示するアプリケーションです。

## 技術スタック
- バックエンド: Go
- フロントエンド: React(TypeScript)
- 通信: WebSocket

## ディレクトリ構成
- backend/: Go バックエンド
- frontend/: React フロントエンド
- docs/: ドキュメント

## コマンド
- make test: 全テスト実行
- make lint: 全体リンティング
- make lint-backend / make test-backend: バックエンド用
- make lint-frontend / make test-frontend: フロントエンド用

## コーディング規約
- クリーンアーキテクチャを厳守し、各層の責務を明確に分離
- ドメイン層はビジネスロジック中心、インフラ/プレゼンテーション層から独立
- 依存関係逆転・インターフェース定義・DI を活用
- DI には `wire` を使用
- TDD を意識し、まずテストを書く
- コメントは日本語で記載
- 関数名から動作を予測できる場合はコメントは不要。複雑なロジックや意図が難解な場合は詳細なコメントを追加
- エラーハンドリング・可読性・テストカバレッジ重視
- エラー時は `err` をそのまま返すのではなく、`fmt.Errorf` でラップして詳細な情報を付加
- 関数・変数名は意味のあるものを使用
- コードの重複を避け、DRY原則を遵守
- ドキュメントは英語と日本語の両方で書き、日本語のドキュメントは `_ja` 接尾辞を付ける
- log.Println や fmt.Println は使用せず、`github.com/nasshu2916/dmx_viewer/pkg/logger` パッケージを利用
  - logger パッケージの関数は `msg string, fields ...interface{}` の形式で、`msg` にメッセージを、`fields` にキーと値のペアを渡す

## ワークフロー
- 変更後は make lint でスタイルチェック
- テストを必ず実行し、通過を確認
- 新機能追加時はユニットテスト必須

## 注意事項
- 指示がない限り git 操作禁止
- go.mod の module 名は `github.com/nasshu2916/dmx_viewer` とする

このルール・ガイドラインに従い、プロジェクトの品質と保守性を高めてください。

