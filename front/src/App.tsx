import { useEffect, useState } from 'react'
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Cell } from 'recharts'
import './App.css'

type Spendings = Record<string, number>

interface ChartData {
  name: string
  hours: number
}

const COLORS = ['#8884d8', '#82ca9d', '#ffc658', '#ff7300', '#0088fe', '#00C49F']

function App() {
  const [data, setData] = useState<ChartData[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [startDate, setStartDate] = useState('2024-01-01')
  const [endDate, setEndDate] = useState('2024-01-07')

  const fetchSpendings = async () => {
    setLoading(true)
    setError(null)
    try {
      const response = await fetch(`/api/calendar/spending?start=${startDate}&end=${endDate}`)
      if (!response.ok) {
        throw new Error(`HTTP ${response.status}`)
      }
      const spendings: Spendings = await response.json()
      const chartData = Object.entries(spendings).map(([name, hours]) => ({
        name,
        hours,
      }))
      setData(chartData)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch data')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchSpendings()
  }, [])

  return (
    <div className="container">
      <h1>Energy Journal - Time Spending</h1>

      <div className="controls">
        <label>
          Start:
          <input
            type="date"
            value={startDate}
            onChange={(e) => setStartDate(e.target.value)}
          />
        </label>
        <label>
          End:
          <input
            type="date"
            value={endDate}
            onChange={(e) => setEndDate(e.target.value)}
          />
        </label>
        <button onClick={fetchSpendings}>Refresh</button>
      </div>

      {loading && <p className="status">Loading...</p>}
      {error && <p className="status error">Error: {error}</p>}

      {!loading && !error && data.length > 0 && (
        <div className="chart-container">
          <ResponsiveContainer width="100%" height={400}>
            <BarChart data={data} margin={{ top: 20, right: 30, left: 20, bottom: 5 }}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" />
              <YAxis label={{ value: 'Hours', angle: -90, position: 'insideLeft' }} />
              <Tooltip formatter={(value) => [`${Number(value).toFixed(2)} hrs`, 'Time Spent']} />
              <Bar dataKey="hours" radius={[4, 4, 0, 0]}>
                {data.map((_, index) => (
                  <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                ))}
              </Bar>
            </BarChart>
          </ResponsiveContainer>
        </div>
      )}
    </div>
  )
}

export default App
