package repository

// WebSocketRepository WebSocketメッセージの送信を抽象化するリポジトリインターフェース
type WebSocketRepository interface {
	// 特定のトピックにメッセージをブロードキャストする
	BroadcastToTopic(topic string, message []byte) error

	// 全クライアントにメッセージをブロードキャストする
	BroadcastToAll(message []byte) error
}
