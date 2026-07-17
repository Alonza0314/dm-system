import { useMemo } from 'react'
import type { Category } from '../../api'
import panelStyles from '../../styles/dashboard-panel.module.css'
import styles from './category-device-bars.module.css'

interface CategoryDeviceBarsProps {
  categories: Category[]
}

interface Row {
  name: string
  idle: number
  using: number
  total: number
}

export default function CategoryDeviceBars({ categories }: CategoryDeviceBarsProps) {
  const rows = useMemo<Row[]>(() => {
    return categories
      .map((category) => {
        const idle = category.idle_device ?? 0
        const using = category.using_device ?? 0
        return { name: category.name ?? '', idle, using, total: idle + using }
      })
      .sort((a, b) => b.total - a.total)
  }, [categories])

  const maxTotal = Math.max(1, ...rows.map((row) => row.total))

  return (
    <div className={styles.wrapper}>
      <p className={styles.heading}>By Category</p>

      {rows.length === 0 ? (
        <p className={panelStyles.tableEmpty}>No categories yet.</p>
      ) : (
        <div className={styles.barList}>
          {rows.map((row) => (
            <div key={row.name} className={styles.barRow}>
              <span className={styles.barLabel} title={row.name}>{row.name}</span>
              <div className={styles.barTrack}>
                <div
                  className={styles.barSegmentIdle}
                  style={{ width: `${(row.idle / maxTotal) * 100}%` }}
                />
                <div
                  className={styles.barSegmentUsing}
                  style={{ width: `${(row.using / maxTotal) * 100}%` }}
                />
              </div>
              <span className={styles.barValue}>{row.total}</span>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
