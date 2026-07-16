import { useEffect, useState, type FormEvent } from 'react'
import Button from '../../components/button/button'
import NotificationContainer from '../../components/notifications/NotificationContainer'
import { useNotifications } from '../../hooks/useNotifications'
import { accountApi, UNAUTHORIZED_MESSAGE_KEY } from '../../apiClient'
import { getErrorMessage } from '../../utils/getErrorMessage'
import { useNavigate } from 'react-router-dom'
import styles from './login-page.module.css'

export default function LoginPage() {
  const navigate = useNavigate()
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [isLoading, setIsLoading] = useState(false)

  const { errors, successes, addError, addSuccess, removeNotification } = useNotifications()

  useEffect(() => {
    const pendingMessage = sessionStorage.getItem(UNAUTHORIZED_MESSAGE_KEY)
    if (pendingMessage) {
      sessionStorage.removeItem(UNAUTHORIZED_MESSAGE_KEY)
      addError(pendingMessage)
    }
  }, [addError])

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setIsLoading(true)

    try {
      const response = await accountApi.login({ username, password })
      const token = response.data.token || ''
      localStorage.setItem('token', token)
      addSuccess(response.data.message || 'Login successful')
      navigate('/', { replace: true })
    } catch (error: unknown) {
      addError(getErrorMessage(error, 'Login failed'))
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className={styles.page}>
      <NotificationContainer
        errors={errors}
        successes={successes}
        onClose={removeNotification}
      />

      <section className={styles.brandTile}>
        <p className={styles.brandKicker}>Device Management</p>
        <h1 className={styles.brandTitle}>DM System</h1>
        <p className={styles.brandSubtitle}>Sign in to continue to your DM System.</p>
      </section>

      <section className={styles.formTile}>
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
            <label className={styles.label} htmlFor="password">Password</label>
            <input
              id="password"
              type="password"
              className={styles.input}
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              autoComplete="current-password"
              required
            />
          </div>

          <div className={styles.actionRow}>
            <Button type="submit" disabled={isLoading}>
              {isLoading ? 'Signing in...' : 'Sign In'}
            </Button>
          </div>
        </form>
      </section>
    </div>
  )
}
