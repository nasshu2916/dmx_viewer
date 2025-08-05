import { LineChart, Line, XAxis, YAxis, Tooltip, CartesianGrid, ResponsiveContainer } from 'recharts'
import type { DmxHistoryPoint } from '@/hooks/useDmxHistory'

interface DmxHistoryChartProps {
  history: DmxHistoryPoint[]
  maxLength: number
}

const DmxHistoryChart = ({ history, maxLength = 100 }: DmxHistoryChartProps) => {
  // 横軸を maxLength 点で固定
  // 横軸を index で固定
  const chartData: { value?: number; index: number }[] = Array.from({ length: maxLength }, (_, i) => ({
    value: undefined,
    index: i,
  }))

  // historyを右詰めでchartDataに埋める
  for (let i = 0; i < history.length && i < maxLength; i++) {
    chartData[maxLength - history.length + i].value = history[i].value
  }

  return (
    <div className="h-32 w-full">
      <ResponsiveContainer height="100%" width="100%">
        <LineChart data={chartData} margin={{ top: 8, right: 8, left: 8, bottom: 8 }}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="index" domain={[0, maxLength - 1]} tickFormatter={t => `${maxLength - t}`} type="number" />
          <YAxis domain={[0, 255]} ticks={[0, 255]} />
          <Tooltip labelFormatter={t => new Date(Number(t)).toLocaleTimeString()} />
          <Line dataKey="value" dot={false} isAnimationActive={false} stroke="#38bdf8" type="monotone" />
        </LineChart>
      </ResponsiveContainer>
    </div>
  )
}

export default DmxHistoryChart
