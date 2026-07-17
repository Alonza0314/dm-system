import { useCallback, useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import type { Device } from '../../api'
import { deviceApi } from '../../apiClient'
import { getErrorMessage } from '../../utils/getErrorMessage'
import panelStyles from '../../styles/dashboard-panel.module.css'
import styles from './device-detail-page.module.css'

const FIELD_LABELS: Array<{ key: keyof Device; label: string }> = [
  { key: 'category', label: 'Category' },
  { key: 'owner', label: 'Owner' },
  { key: 'user', label: 'User' },
  { key: 'note', label: 'Note' },
  { key: 'id', label: 'ID' },
]

export default function DeviceDetailPage() {
  const { categoryName = '', deviceName = '' } = useParams()

  const [device, setDevice] = useState<Device | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

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

  useEffect(() => {
    fetchDevice()
    const interval = setInterval(() => fetchDevice(true), 10_000)
    return () => clearInterval(interval)
  }, [fetchDevice])

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
                <p className={styles.fieldValue}>{device[key] || '—'}</p>
              </div>
            ))}
          </div>
        </>
      )}
    </section>
  )
}
