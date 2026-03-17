import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { marketAPI } from '../services/api'
import { useAuth } from '../context/AuthContext'
import {
  AreaChart, Area, XAxis, YAxis, CartesianGrid,
  Tooltip, ResponsiveContainer
} from 'recharts'
import { ArrowLeft, TrendingUp, TrendingDown, LogOut, Loader2 } from 'lucide-react'

function CustomTooltip({ active, payload, label }) {
  if (!active || !payload?.length) return null
  return (
    <div className="bg-gray-900 border border-gray-700 rounded-lg px-4 py-3 text-sm">
      <p className="text-gray-400 mb-1">{label}</p>
      <p className="text-white font-semibold">Close: {payload[0]?.value?.toFixed(2)}</p>
    </div>
  )
}

export default function StockPage() {
  const { symbol }  = useParams()
  const navigate    = useNavigate()
  const { user, logout } = useAuth()

  const [quote, setQuote]     = useState(null)
  const [history, setHistory] = useState([])
  const [period, setPeriod]   = useState('30d')
  const [loading, setLoading] = useState(true)
  const [error, setError]     = useState('')

  useEffect(() => {
    const load = async () => {
      setLoading(true)
      setError('')
      try {
        const [quoteRes, histRes] = await Promise.all([
          marketAPI.quote(symbol),
          marketAPI.history(symbol, period),
        ])
        setQuote(quoteRes.data)
        const formatted = (histRes.data.data || []).map(d => ({
          ...d,
          date:  new Date(d.trade_date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' }),
          close: parseFloat((d.close || 0).toFixed(2)),
          open:  parseFloat((d.open  || 0).toFixed(2)),
          high:  parseFloat((d.high  || 0).toFixed(2)),
          low:   parseFloat((d.low   || 0).toFixed(2)),
        }))
        setHistory(formatted)
      } catch (err) {
        setError(err.response?.data?.error || 'Failed to load data')
      } finally {
        setLoading(false)
      }
    }
    load()
  }, [symbol, period])

  const isPositive   = (quote?.change ?? 0) >= 0
  const currencySign = quote?.currency === 'INR' ? '₹' : '$'

  return (
    <div className="min-h-screen bg-gray-950 text-white">
      {/* Navbar */}
      <nav className="bg-gray-900 border-b border-gray-800 px-6 py-4 flex items-center justify-between">
        <div className="flex items-center gap-4">
          <button onClick={() => navigate('/')} className="text-gray-400 hover:text-white transition p-1 rounded-lg hover:bg-gray-800">
            <ArrowLeft size={20} />
          </button>
          <div className="flex items-center gap-2">
            <TrendingUp className="text-blue-500" size={22} />
            <span className="font-bold">TradePulse</span>
          </div>
        </div>
        <div className="flex items-center gap-4">
          <span className="text-gray-400 text-sm">{user?.name}</span>
          <button onClick={logout} className="flex items-center gap-1.5 text-gray-400 hover:text-white text-sm transition">
            <LogOut size={16} /> Logout
          </button>
        </div>
      </nav>

      <main className="max-w-5xl mx-auto px-6 py-8">
        {loading ? (
          <div className="flex flex-col items-center justify-center py-24 gap-4">
            <Loader2 className="animate-spin text-blue-500" size={40} />
            <p className="text-gray-500">Fetching data for {symbol}...</p>
          </div>
        ) : error ? (
          <div className="bg-red-500/10 border border-red-500/30 text-red-400 rounded-xl p-8 text-center">
            <p className="font-semibold mb-1">Failed to load {symbol}</p>
            <p className="text-sm">{error}</p>
            <button onClick={() => navigate('/')} className="mt-4 text-blue-400 hover:underline text-sm">← Back to search</button>
          </div>
        ) : (
          <>
            {/* Quote Card */}
            {quote && (
              <div className="bg-gray-900 border border-gray-800 rounded-2xl p-6 mb-6">
                <div className="flex items-start justify-between gap-4">
                  <div>
                    <h1 className="text-2xl font-bold">{quote.symbol}</h1>
                    <p className="text-gray-400 text-sm mt-0.5">{quote.company_name}</p>
                    <p className="text-xs text-gray-600 mt-1">{quote.exchange} · {quote.currency}</p>
                  </div>
                  <div className="text-right">
                    <p className="text-3xl font-bold tabular-nums">
                      {currencySign}{quote.current_price?.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                    </p>
                    <div className={`flex items-center justify-end gap-1 mt-1 font-medium ${isPositive ? 'text-green-400' : 'text-red-400'}`}>
                      {isPositive ? <TrendingUp size={16} /> : <TrendingDown size={16} />}
                      <span>{isPositive ? '+' : ''}{quote.change?.toFixed(2)} ({isPositive ? '+' : ''}{quote.change_percent?.toFixed(2)}%)</span>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {/* Period Tabs */}
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold">Price Chart</h2>
              <div className="flex gap-2 bg-gray-900 border border-gray-800 rounded-xl p-1">
                {[{ value: '1d', label: '1 Day' }, { value: '30d', label: '30 Days' }].map(({ value, label }) => (
                  <button key={value} onClick={() => setPeriod(value)}
                    className={`px-4 py-1.5 rounded-lg text-sm font-medium transition ${
                      period === value ? 'bg-blue-600 text-white' : 'text-gray-400 hover:text-white hover:bg-gray-800'
                    }`}>
                    {label}
                  </button>
                ))}
              </div>
            </div>

            {/* Chart */}
            <div className="bg-gray-900 border border-gray-800 rounded-2xl p-6 mb-6">
              {history.length === 0 ? (
                <p className="text-center text-gray-500 py-14">No chart data available</p>
              ) : (
                <ResponsiveContainer width="100%" height={320}>
                  <AreaChart data={history} margin={{ top: 5, right: 10, left: 0, bottom: 5 }}>
                    <defs>
                      <linearGradient id="gradClose" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="5%"  stopColor="#3b82f6" stopOpacity={0.25} />
                        <stop offset="95%" stopColor="#3b82f6" stopOpacity={0} />
                      </linearGradient>
                    </defs>
                    <CartesianGrid strokeDasharray="3 3" stroke="#1f2937" vertical={false} />
                    <XAxis dataKey="date" tick={{ fill: '#6b7280', fontSize: 12 }} tickLine={false} axisLine={false} interval="preserveStartEnd" />
                    <YAxis domain={['auto', 'auto']} tick={{ fill: '#6b7280', fontSize: 12 }} tickLine={false} axisLine={false} width={70}
                      tickFormatter={v => `${currencySign}${v.toFixed(0)}`} />
                    <Tooltip content={<CustomTooltip />} />
                    <Area type="monotone" dataKey="close" stroke="#3b82f6" strokeWidth={2} fill="url(#gradClose)"
                      dot={false} activeDot={{ r: 5, fill: '#3b82f6', strokeWidth: 0 }} />
                  </AreaChart>
                </ResponsiveContainer>
              )}
            </div>

            {/* Table */}
            <div className="bg-gray-900 border border-gray-800 rounded-2xl overflow-hidden">
              <div className="px-6 py-4 border-b border-gray-800 flex items-center justify-between">
                <h3 className="font-semibold">Historical Data</h3>
                <span className="text-gray-500 text-sm">{history.length} records</span>
              </div>
              <div className="overflow-x-auto">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="text-gray-500 text-xs uppercase border-b border-gray-800">
                      {['Date', 'Open', 'High', 'Low', 'Close', 'Volume'].map(h => (
                        <th key={h} className="px-5 py-3 text-right first:text-left font-medium">{h}</th>
                      ))}
                    </tr>
                  </thead>
                  <tbody>
                    {[...history].reverse().map((row, i) => (
                      <tr key={i} className="border-b border-gray-800/40 hover:bg-gray-800/30 transition">
                        <td className="px-5 py-3 text-gray-300">{row.date}</td>
                        <td className="px-5 py-3 text-right tabular-nums">{row.open?.toFixed(2)}</td>
                        <td className="px-5 py-3 text-right tabular-nums text-green-400">{row.high?.toFixed(2)}</td>
                        <td className="px-5 py-3 text-right tabular-nums text-red-400">{row.low?.toFixed(2)}</td>
                        <td className="px-5 py-3 text-right tabular-nums font-medium">{row.close?.toFixed(2)}</td>
                        <td className="px-5 py-3 text-right tabular-nums text-gray-400">
                          {row.volume ? row.volume.toLocaleString() : '—'}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          </>
        )}
      </main>
    </div>
  )
}