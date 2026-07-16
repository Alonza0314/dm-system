import { useState, type FormEvent } from 'react'
import { useParams } from 'react-router-dom'
import Button from '../../components/button/button'
import NotificationContainer from '../../components/notifications/NotificationContainer'
import { useNotifications } from '../../hooks/useNotifications'
import styles from './qrcode-page.module.css'

type DeviceAction = 'borrow' | 'return'

export default function QrCodePage() {
  const { cate = '', dev = '' } = useParams()
  const { errors, successes, addError, addSuccess, removeNotification } = useNotifications()

  const [username, setUsername] = useState('')
  const [action, setAction] = useState<DeviceAction | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)

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

    setIsSubmitting(true)
    try {
      // TODO: wire this up once the no-auth device-record backend endpoint exists.
      // Expected to POST { category: cate, name: dev, username: trimmedUsername, action } here.
      addSuccess(`${action === 'borrow' ? 'Borrow' : 'Return'} recorded for "${trimmedUsername}"`)
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <div className={styles.page}>
      <NotificationContainer errors={errors} successes={successes} onClose={removeNotification} />

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
    </div>
  )
}
