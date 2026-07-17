import { useMemo, useState, type MouseEvent } from 'react'
import type { Category } from '../../api'
import panelStyles from '../../styles/dashboard-panel.module.css'
import styles from './status-overview.module.css'

interface StatusOverviewProps {
  categories: Category[]
}

interface Breakdown {
  name: string
  count: number
}

type HoveredStatus = 'idle' | 'using' | null

const RADIUS = 45
const STROKE = 16
const STROKE_ACTIVE = 20
const CIRCUMFERENCE = 2 * Math.PI * RADIUS
const GAP = 4

function sortedBreakdown(categories: Category[], key: 'idle_device' | 'using_device'): Breakdown[] {
  return categories
    .map((category) => ({ name: category.name ?? '', count: category[key] ?? 0 }))
    .filter((entry) => entry.count > 0)
    .sort((a, b) => b.count - a.count)
}

export default function StatusOverview({ categories }: StatusOverviewProps) {
  const [hovered, setHovered] = useState<HoveredStatus>(null)
  const [tooltipPos, setTooltipPos] = useState({ x: 0, y: 0 })

  const totalIdle = useMemo(
    () => categories.reduce((sum, category) => sum + (category.idle_device ?? 0), 0),
    [categories],
  )
  const totalUsing = useMemo(
    () => categories.reduce((sum, category) => sum + (category.using_device ?? 0), 0),
    [categories],
  )
  const total = totalIdle + totalUsing

  const idleBreakdown = useMemo(() => sortedBreakdown(categories, 'idle_device'), [categories])
  const usingBreakdown = useMemo(() => sortedBreakdown(categories, 'using_device'), [categories])

  const idleLen = total > 0 ? (totalIdle / total) * CIRCUMFERENCE : 0
  const usingLen = total > 0 ? (totalUsing / total) * CIRCUMFERENCE : 0

  function handleMove(event: MouseEvent<Element>, status: HoveredStatus) {
    const container = event.currentTarget.closest(`.${styles.chartRow}`) as HTMLElement | null
    if (!container) return
    const rect = container.getBoundingClientRect()
    setTooltipPos({ x: event.clientX - rect.left, y: event.clientY - rect.top })
    setHovered(status)
  }

  function handleLeave() {
    setHovered(null)
  }

  const isIdleActive = hovered === 'idle'
  const isUsingActive = hovered === 'using'
  const activeBreakdown = isIdleActive ? idleBreakdown : isUsingActive ? usingBreakdown : null

  return (
    <div className={panelStyles.tile}>
      <div className={panelStyles.tileHeader}>
        <div>
          <p className={panelStyles.tileTag}>Devices</p>
          <h3>Quick View</h3>
        </div>
      </div>

      <div className={styles.chartRow}>
        <svg viewBox="0 0 120 120" className={styles.donut}>
          {total === 0 ? (
            <circle cx="60" cy="60" r={RADIUS} fill="none" stroke="var(--color-hairline)" strokeWidth={STROKE} />
          ) : (
            <>
              <circle
                cx="60"
                cy="60"
                r={RADIUS}
                fill="none"
                stroke="var(--color-status-idle)"
                strokeWidth={isIdleActive ? STROKE_ACTIVE : STROKE}
                strokeLinecap="round"
                strokeDasharray={`${Math.max(idleLen - GAP, 0)} ${CIRCUMFERENCE}`}
                strokeDashoffset={0}
                transform="rotate(-90 60 60)"
                onMouseMove={(event) => handleMove(event, 'idle')}
                onMouseLeave={handleLeave}
                className={`${styles.arc} ${isUsingActive ? styles.arcDimmed : ''}`}
              />
              <circle
                cx="60"
                cy="60"
                r={RADIUS}
                fill="none"
                stroke="var(--color-status-using)"
                strokeWidth={isUsingActive ? STROKE_ACTIVE : STROKE}
                strokeLinecap="round"
                strokeDasharray={`${Math.max(usingLen - GAP, 0)} ${CIRCUMFERENCE}`}
                strokeDashoffset={-idleLen}
                transform="rotate(-90 60 60)"
                onMouseMove={(event) => handleMove(event, 'using')}
                onMouseLeave={handleLeave}
                className={`${styles.arc} ${isUsingActive ? styles.arcActive : ''} ${
                  isIdleActive ? styles.arcDimmed : ''
                }`}
              />
            </>
          )}
          <text x="60" y="56" textAnchor="middle" className={styles.totalValue}>
            {isIdleActive ? totalIdle : isUsingActive ? totalUsing : total}
          </text>
          <text x="60" y="72" textAnchor="middle" className={styles.totalLabel}>
            {isIdleActive ? 'idle' : isUsingActive ? 'using' : 'devices'}
          </text>
        </svg>

        <div className={styles.legend}>
          <div
            className={`${styles.legendRow} ${isIdleActive ? styles.legendRowActiveIdle : ''}`}
            onMouseMove={(event) => handleMove(event, 'idle')}
            onMouseLeave={handleLeave}
          >
            <span className={`${styles.swatch} ${styles.swatchIdle}`} />
            <span className={styles.legendLabel}>Idle</span>
            <span className={styles.legendValue}>{totalIdle}</span>
          </div>
          <div
            className={`${styles.legendRow} ${isUsingActive ? styles.legendRowActiveUsing : ''}`}
            onMouseMove={(event) => handleMove(event, 'using')}
            onMouseLeave={handleLeave}
          >
            <span className={`${styles.swatch} ${styles.swatchUsing}`} />
            <span className={styles.legendLabel}>Using</span>
            <span className={styles.legendValue}>{totalUsing}</span>
          </div>
        </div>

        {activeBreakdown && (
          <div
            className={styles.tooltip}
            style={{ left: tooltipPos.x + 16, top: tooltipPos.y + 16 }}
          >
            <p className={styles.tooltipTitle}>
              {hovered === 'idle' ? 'Idle by category' : 'Using by category'}
            </p>
            {activeBreakdown.length === 0 ? (
              <p className={styles.tooltipEmpty}>None</p>
            ) : (
              <ul className={styles.tooltipList}>
                {activeBreakdown.map((entry) => (
                  <li key={entry.name} className={styles.tooltipRow}>
                    <span>{entry.name}</span>
                    <span>{entry.count}</span>
                  </li>
                ))}
              </ul>
            )}
          </div>
        )}
      </div>
    </div>
  )
}
