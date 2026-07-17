import { useState, type FormEvent } from 'react'
import { useParams } from 'react-router-dom'
import Button from '../../components/button/button'
import NotificationContainer from '../../components/notifications/NotificationContainer'
import { useNotifications } from '../../hooks/useNotifications'
import { qrcodeApi } from '../../apiClient'
import { getErrorMessage } from '../../utils/getErrorMessage'
import styles from './qrcode-page.module.css'

type DeviceAction = 'borrow' | 'return'
type ResultState = { type: 'success' } | { type: 'error'; message: string }

const USERNAME_STORAGE_KEY = 'dm-system:qrcode-username'

export default function QrCodePage() {
  const { cate = '', dev = '' } = useParams()
  const { errors, successes, addError, removeNotification } = useNotifications()

  const [username, setUsername] = useState(() => localStorage.getItem(USERNAME_STORAGE_KEY) || '')
  const [action, setAction] = useState<DeviceAction | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [result, setResult] = useState<ResultState | null>(null)

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const trimmedUsername = username.trim()
    if (!trimmedUsername) {
      addError('Username is required')
      return
    }
    if (!action) {
      addError('Please choose Borrow or Return')
      return
    }

    localStorage.setItem(USERNAME_STORAGE_KEY, trimmedUsername)

    setIsSubmitting(true)
    try {
      if (action === 'borrow') {
        await qrcodeApi.borrowDevice(cate, dev, { user: trimmedUsername })
      } else {
        await qrcodeApi.returnDevice(cate, dev, { user: trimmedUsername })
      }
      setResult({ type: 'success' })
    } catch (err) {
      const status = (err as { response?: { status?: number } })?.response?.status
      let message = getErrorMessage(err, `Failed to ${action === 'borrow' ? 'borrow' : 'return'} device`)
      if (status === 500) {
        message = `${message} Please contact IT admin.`
      }
      setResult({ type: 'error', message })
    } finally {
      setIsSubmitting(false)
    }
  }

  function handleReset() {
    setResult(null)
    setAction(null)
  }

  return (
    <div className={styles.page}>
      <NotificationContainer errors={errors} successes={successes} onClose={removeNotification} />

      {result ? (
        <main
          className={`${styles.card} ${
            result.type === 'success' ? styles.cardSuccess : styles.cardError
          }`}
        >
          <p className={styles.resultIcon}>{result.type === 'success' ? '✓' : '✕'}</p>
          <h1 className={styles.resultTitle}>
            {result.type === 'success' ? 'Success' : 'Something went wrong'}
          </h1>
          <p className={styles.resultMessage}>
            {result.type === 'success'
              ? `${action === 'borrow' ? 'Borrow' : 'Return'} recorded for "${username.trim()}".`
              : result.message}
          </p>
          <div className={styles.actionRow}>
            <Button onClick={handleReset}>{result.type === 'success' ? 'Go Back' : 'Try Again'}</Button>
          </div>
        </main>
      ) : (
        <main className={styles.card}>
          <div className={styles.headerBlock}>
            <p className={styles.kicker}>Device Check-in / Check-out</p>
            <h1 className={styles.title}>{dev}</h1>
            <p className={styles.subtitle}>Category: {cate}</p>
          </div>

          <form className={styles.form} onSubmit={handleSubmit}>
            <div>
              <label className={styles.label} htmlFor="username">Username</label>
              <input
                id="username"
                className={styles.input}
                value={username}
                onChange={(event) => setUsername(event.target.value)}
                autoComplete="username"
                required
              />
            </div>

            <div>
              <p className={styles.label}>Action</p>
              <div className={styles.actionToggle} role="radiogroup" aria-label="Action">
                <label className={`${styles.actionOption} ${action === 'borrow' ? styles.actionOptionSelected : ''}`}>
                  <input
                    type="radio"
                    name="action"
                    value="borrow"
                    checked={action === 'borrow'}
                    onChange={() => setAction('borrow')}
                  />
                  Borrow
                </label>
                <label className={`${styles.actionOption} ${action === 'return' ? styles.actionOptionSelected : ''}`}>
                  <input
                    type="radio"
                    name="action"
                    value="return"
                    checked={action === 'return'}
                    onChange={() => setAction('return')}
                  />
                  Return
                </label>
              </div>
            </div>

            <div className={styles.actionRow}>
              <Button type="submit" disabled={isSubmitting}>
                {isSubmitting ? 'Submitting...' : 'Submit'}
              </Button>
            </div>
          </form>
        </main>
      )}
    </div>
  )
}
