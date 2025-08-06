/**
 * columns（2のN乗）を計算するユーティリティ関数
 * @param containerWidth - コンテナの幅(px)
 * @param cellWidth - セルの最小幅(px)
 * @param minColumns - 最小カラム数（デフォルト1）
 * @param maxColumns - 最大カラム数（デフォルト32）
 * @returns 2のN乗のカラム数（minColumns <= columns <= maxColumns）
 */
export function calcColumns(containerWidth: number, cellWidth = 24, minColumns = 1, maxColumns = 32): number {
  let cols = Math.max(minColumns, Math.min(maxColumns, Math.floor(containerWidth / cellWidth)))
  // 2のN乗に切り捨て
  cols = Math.pow(2, Math.floor(Math.log2(cols)))
  if (cols < minColumns) cols = minColumns
  if (cols > maxColumns) cols = maxColumns
  return cols
}
