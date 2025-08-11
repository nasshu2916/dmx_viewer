# DMX Viewer

## 概要

DMX Viewer は、ArtNet 信号をリアルタイムで可視化する Web ページを提供します。

![dmx_viewer](https://github.com/user-attachments/assets/da3a52e7-0712-4b6e-8bbc-7249f7bb9c4e)

## 機能

- **リアルタイムでDMX信号を表示**  
  受信したArtNetの信号を、表やグラフで表示します。
- **ノードごとに受信したユニバースを表示**  
  各ノードごとに受信したユニバースの状態を表示します。
- **履歴グラフで変化を確認**  
  DMX値の変化を時系列で表示します。

## アーキテクチャ

- **フロントエンド**
  - フレームワーク: React (TypeScript)
  - バックエンド: Go
  - 通信: WebSocket
