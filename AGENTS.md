# DMX Viewer ガイドライン

---

## 1. プロジェクト概要・技術スタック

このプロジェクトはバックエンドで受信した ArtNet の情報を WebSocket 経由でフロントエンドに送信し、リアルタイムで DMX データを表示するアプリケーションです。

- バックエンド: Go
- フロントエンド: React(TypeScript)
- 通信: WebSocket

---

## 2. ディレクトリ構成

- backend/: Go バックエンド
- frontend/: React フロントエンド
- docs/: ドキュメント

---

## 3. ビルド・テスト・リントコマンド

- 全テスト: `make test`
- バックエンドテスト: `make test-backend` または `cd backend && go test ./...`
- フロントエンドテスト: `npm run test --prefix frontend` または `cd frontend && npm test`
- Go パッケージ単体テスト: `cd backend && go test ./internal/usecase`
- TypeScript ファイル単体テスト: `cd frontend && npm test -- TimeDisplay.test.tsx`
- 全体リンティング: `make lint`
- バックエンド lint: `cd backend && go fmt ./... && go vet ./...`
- フロントエンド lint: `cd frontend && npm run lint && npm run type-check`

---

## 4. アーキテクチャ・設計原則

- クリーンアーキテクチャを厳守し、各層の責務を明確に分離
- ドメイン層はビジネスロジック中心、インフラ/プレゼンテーション層から独立
- 依存関係逆転・インターフェース定義・DI を活用
- DI には `wire` を使用
- TDD を意識し、まずテストを書く
- コードの重複を避け、DRY原則を遵守

---

## 5. Go バックエンド コーディング規約

- 依存性注入は Google Wire を利用
- パッケージ構成: domain/model, usecase, interface/handler, infrastructure
- エラーハンドリング: panic禁止、エラーは return で返す
- エラー時は `err` をそのまま返さず、`fmt.Errorf` でラップして詳細情報を付加
- 命名: エクスポートは PascalCase、非公開は camelCase
- インポート順: 標準ライブラリ→サードパーティ→内部
- 関数・変数名は意味のあるものを使用
- コメントは日本語で記載。関数名から動作が明確な場合は不要、複雑な場合は詳細なコメントを追加
- log.Println や fmt.Println は使用せず、`pkg/logger` パッケージを利用
  - logger関数は `msg string, fields ...interface{}` 形式

---

## 6. TypeScript フロントエンド コーディング規約

- React関数コンポーネント＋hooks
- 厳格なTypeScript: `any`禁止、明示的な型指定
- Prettier: single quotes, no semicolons, 2 spaces, 120文字幅
- Propsはshorthand→callbackの順で並べる
- TailwindCSS（`dmx-` prefixカスタム）
- エラーハンドリング: try-catch＋console.warn（致命的でない場合）
- パフォーマンスを意識し、不要な再レンダリングを極力抑える
  - useMemo, useCallback, React.memo などを適切に活用する

---

## 7. ロギング・ドキュメント

- loggerは `github.com/nasshu2916/dmx_viewer/pkg/logger` または `logger.ts` を使用
  - `logger` 関数は `msg string, fields ...interface{}` の形式で、`msg` にメッセージ、`fields` にキーと値のペアを渡す
  - ログレベルは `info`, `warn`, `error` を使用
- ドキュメントは英語・日本語両方で作成。日本語は `_ja` サフィックス

---

## 8. ワークフロー

- 変更後は `make lint` でスタイルチェック
- テストを必ず実行し、通過を確認
- 新機能追加時はユニットテスト必須

---

## 9. 注意事項

- 指示がない限り git 操作禁止
- go.mod の module 名は `github.com/nasshu2916/dmx_viewer` とする

---

このルール・ガイドラインに従い、プロジェクトの品質と保守性を高めてください。
