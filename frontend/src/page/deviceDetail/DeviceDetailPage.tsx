import { useCallback, useEffect, useMemo, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import type { Device, HistoryEntry } from '../../api'
import { deviceApi, historyApi } from '../../apiClient'
import { getErrorMessage } from '../../utils/getErrorMessage'
import Button from '../../components/button/button'
import panelStyles from '../../styles/dashboard-panel.module.css'
import styles from './device-detail-page.module.css'

const FIELD_LABELS: Array<{ key: keyof Device; label: string }> = [
  { key: 'category', label: 'Category' },
  { key: 'owner', label: 'Owner' },
  { key: 'user', label: 'User' },
  { key: 'lastBorrow', label: 'Last Borrow' },
  { key: 'lastReturn', label: 'Last Return' },
  { key: 'note', label: 'Note' },
  { key: 'id', label: 'ID' },
]

const TIMESTAMP_FIELDS = new Set<keyof Device>(['lastBorrow', 'lastReturn'])

const HISTORY_PAGE_SIZE = 10

function formatFieldValue(key: keyof Device, value: Device[keyof Device]) {
  if (!value) return '—'
  if (TIMESTAMP_FIELDS.has(key)) return new Date(value as string).toLocaleString()
  return value
}

function formatTimestamp(value?: string | null) {
  if (!value) return '—'
  return new Date(value).toLocaleString()
}

export default function DeviceDetailPage() {
  const { categoryName = '', deviceName = '' } = useParams()

  const [device, setDevice] = useState<Device | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const [history, setHistory] = useState<HistoryEntry[]>([])
  const [isHistoryLoading, setIsHistoryLoading] = useState(true)
  const [historyError, setHistoryError] = useState<string | null>(null)
  const [historyPage, setHistoryPage] = useState(1)

  const fetchDevice = useCallback((silent = false) => {
    if (!silent) setIsLoading(true)
    return deviceApi.getDevice(categoryName, deviceName)
      .then((response) => {
        setDevice(response.data)
        setError(null)
      })
      .catch((err) => {
        setError(getErrorMessage(err, 'Failed to load device'))
      })
      .finally(() => {
        if (!silent) setIsLoading(false)
      })
  }, [categoryName, deviceName])

  const fetchHistory = useCallback((silent = false) => {
    if (!silent) setIsHistoryLoading(true)
    return historyApi.getHistory(categoryName, deviceName)
      .then((response) => {
        setHistory(response.data.Histories ?? [])
        setHistoryError(null)
      })
      .catch((err) => {
        setHistoryError(getErrorMessage(err, 'Failed to load history'))
      })
      .finally(() => {
        if (!silent) setIsHistoryLoading(false)
      })
  }, [categoryName, deviceName])

  useEffect(() => {
    fetchDevice()
    const interval = setInterval(() => fetchDevice(true), 10_000)
    return () => clearInterval(interval)
  }, [fetchDevice])

  useEffect(() => {
    setHistoryPage(1)
    fetchHistory()
  }, [fetchHistory])

  const historyPageCount = Math.max(1, Math.ceil(history.length / HISTORY_PAGE_SIZE))

  const visibleHistory = useMemo(() => {
    const start = (historyPage - 1) * HISTORY_PAGE_SIZE
    return history.slice(start, start + HISTORY_PAGE_SIZE)
  }, [history, historyPage])

  return (
    <section className={panelStyles.tile}>
      <div className={panelStyles.tileHeader}>
        <div>
          <p className={panelStyles.tileTag}>Device</p>
          <h3>{deviceName}</h3>
        </div>
        <Link to={`/category/${encodeURIComponent(categoryName)}`} className={styles.backLink}>
          Back to {categoryName}
        </Link>
      </div>

      {error && <p className={panelStyles.tableError}>{error}</p>}
      {isLoading && !device && !error && <p className={styles.loading}>Loading...</p>}

      {device && (
        <>
          <div className={styles.identityCard}>
            <h4 className={styles.identityName}>{device.name}</h4>
            <span
              className={`${panelStyles.statusPill} ${styles.statusPillLarge} ${
                device.status === 'using' ? panelStyles.statusUsing : panelStyles.statusIdle
              }`}
            >
              {device.status}
            </span>
          </div>

          <div className={styles.fieldGrid}>
            {FIELD_LABELS.map(({ key, label }) => (
              <div key={key} className={styles.fieldCard}>
                <p className={styles.fieldLabel}>{label}</p>
                <p className={styles.fieldValue}>{formatFieldValue(key, device[key])}</p>
              </div>
            ))}
          </div>

          <h4 className={styles.historyTitle}>History</h4>

          {historyError && <p className={panelStyles.tableError}>{historyError}</p>}

          <div className={panelStyles.tableWrap}>
            <table className={panelStyles.table}>
              <thead>
                <tr>
                  <th>User</th>
                  <th>Last Borrow</th>
                  <th>Last Return</th>
                </tr>
              </thead>
              <tbody>
                {visibleHistory.map((entry, index) => (
                  <tr key={`${entry.user}-${entry.lastBorrow}-${index}`}>
                    <td>{entry.user || '—'}</td>
                    <td>{formatTimestamp(entry.lastBorrow)}</td>
                    <td>{formatTimestamp(entry.lastReturn)}</td>
                  </tr>
                ))}
                {!isHistoryLoading && history.length === 0 && (
                  <tr>
                    <td colSpan={3} className={panelStyles.tableEmpty}>No history yet.</td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>

          {history.length > 0 && (
            <div className={styles.paginationBar}>
              <Button
                variant="secondary"
                onClick={() => setHistoryPage((page) => Math.max(1, page - 1))}
                disabled={historyPage <= 1}
              >
                Previous
              </Button>
              <span className={styles.paginationLabel}>
                Page {historyPage} of {historyPageCount}
              </span>
              <Button
                variant="secondary"
                onClick={() => setHistoryPage((page) => Math.min(historyPageCount, page + 1))}
                disabled={historyPage >= historyPageCount}
              >
                Next
              </Button>
            </div>
          )}
        </>
      )}
    </section>
  )
}
