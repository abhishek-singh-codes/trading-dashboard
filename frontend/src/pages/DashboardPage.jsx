import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { marketAPI } from '../services/api'
import { useAuth } from '../context/AuthContext'
import { Search, TrendingUp, LogOut, Loader2, BarChart2 } from 'lucide-react'

const QUICK_PICKS = [
  { symbol: 'AAPL',        label: 'Apple' },
  { symbol: 'GOOGL',       label: 'Google' },
  { symbol: 'MSFT',        label: 'Microsoft' },
  { symbol: 'TSLA',        label: 'Tesla' },
  { symbol: 'AMZN',        label: 'Amazon' },
  { symbol: 'RELIANCE.NS', label: 'Reliance' },
  { symbol: 'TCS.NS',      label: 'TCS' },
  { symbol: 'INFY.NS',     label: 'Infosys' },
]

export default function DashboardPage() {
  const [query, setQuery]       = useState('')
  const [results, setResults]   = useState([])
  const [loading, setLoading]   = useState(false)
  const [searched, setSearched] = useState(false)
  const [error, setError]       = useState('')
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  const handleSearch = async (e) => {
    e.preventDefault()
    if (!query.trim()) return
    setLoading(true)
    setSearched(true)
    setError('')
    try {
      const { data } = await marketAPI.search(query.trim())
      setResults(data.results || [])
    } catch (err) {
      setError('Search failed. Please try again.')
      setResults([])
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-gray-950 text-white">
      {/* Navbar */}
      <nav className="bg-gray-900 border-b border-gray-800 px-6 py-4 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <TrendingUp className="text-blue-500" size={24} />
          <span className="font-bold text-lg">TradePulse</span>
        </div>
        <div className="flex items-center gap-4">
          <span className="text-gray-400 text-sm">Hello, {user?.name}</span>
          <button onClick={logout} className="flex items-center gap-1.5 text-gray-400 hover:text-white transition text-sm">
            <LogOut size={16} /> Logout
          </button>
        </div>
      </nav>

      <main className="max-w-4xl mx-auto px-6 py-12">
        <div className="text-center mb-10">
          <BarChart2 className="text-blue-500 mx-auto mb-3" size={36} />
          <h1 className="text-3xl font-bold mb-2">Trading Dashboard</h1>
          <p className="text-gray-400">Search for any stock, ETF, or mutual fund</p>
        </div>

        {/* Search */}
        <form onSubmit={handleSearch} className="flex gap-3 mb-8">
          <div className="relative flex-1">
            <Search className="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400" size={18} />
            <input
              type="text" value={query}
              onChange={e => setQuery(e.target.value)}
              placeholder="Search stocks... (e.g. Apple, RELIANCE)"
              className="w-full bg-gray-800 border border-gray-700 text-white rounded-xl pl-11 pr-4 py-3 focus:outline-none focus:border-blue-500 transition placeholder-gray-500"
            />
          </div>
          <button type="submit" disabled={loading}
            className="bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white px-6 rounded-xl font-medium transition flex items-center gap-2">
            {loading ? <Loader2 className="animate-spin" size={18} /> : <Search size={18} />}
            Search
          </button>
        </form>

        {/* Quick Picks */}
        {!searched && (
          <div>
            <p className="text-gray-500 text-sm mb-3 font-medium uppercase tracking-wide">Popular Stocks</p>
            <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
              {QUICK_PICKS.map(({ symbol, label }) => (
                <button key={symbol} onClick={() => navigate(`/stock/${symbol}`)}
                  className="bg-gray-900 hover:bg-gray-800 border border-gray-800 hover:border-blue-500/50 text-left px-4 py-3 rounded-xl transition group">
                  <p className="font-semibold text-white text-sm group-hover:text-blue-400 transition">{symbol}</p>
                  <p className="text-gray-500 text-xs mt-0.5">{label}</p>
                </button>
              ))}
            </div>
          </div>
        )}

        {/* Results */}
        {searched && (
          <div>
            {loading ? (
              <div className="flex justify-center py-16">
                <Loader2 className="animate-spin text-blue-500" size={36} />
              </div>
            ) : error ? (
              <div className="bg-red-500/10 border border-red-500/30 text-red-400 rounded-xl p-6 text-center">{error}</div>
            ) : results.length === 0 ? (
              <div className="text-center py-12">
                <p className="text-gray-400">No results for &quot;{query}&quot;</p>
                <button onClick={() => { setSearched(false); setQuery('') }}
                  className="text-blue-400 text-sm mt-3 hover:underline">← Back</button>
              </div>
            ) : (
              <>
                <p className="text-gray-500 text-sm mb-3">{results.length} results for &quot;{query}&quot;</p>
                <div className="space-y-2">
                  {results.map((r) => (
                    <button key={r.symbol} onClick={() => navigate(`/stock/${r.symbol}`)}
                      className="w-full text-left bg-gray-900 hover:bg-gray-800 border border-gray-800 hover:border-blue-500/40 rounded-xl px-5 py-4 flex items-center justify-between transition">
                      <div>
                        <span className="font-semibold text-white">{r.symbol}</span>
                        <p className="text-gray-400 text-sm mt-0.5">{r.name || 'N/A'}</p>
                      </div>
                      <div className="text-right">
                        <span className="text-xs bg-gray-700 text-gray-300 px-2 py-1 rounded-md">{r.quote_type}</span>
                        <p className="text-gray-600 text-xs mt-1">{r.exchange}</p>
                      </div>
                    </button>
                  ))}
                </div>
              </>
            )}
          </div>
        )}
      </main>
    </div>
  )
}