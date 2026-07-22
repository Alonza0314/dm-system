import { useState, type FormEvent } from 'react'
import { useNavigate } from 'react-router-dom'
import Button from '../../components/button/button'
import NotificationContainer from '../../components/notifications/NotificationContainer'
import { useNotifications } from '../../hooks/useNotifications'
import { settingApi, PENDING_SUCCESS_MESSAGE_KEY } from '../../apiClient'
import { getErrorMessage } from '../../utils/getErrorMessage'
import panelStyles from '../../styles/dashboard-panel.module.css'
import modalStyles from '../../components/modal/modal.module.css'
import styles from './settings-page.module.css'

export default function SettingsPage() {
  const navigate = useNavigate()
  const { errors, successes, addError, removeNotification } = useNotifications()

  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const trimmedUsername = username.trim()
    if (!trimmedUsername || !password || !confirmPassword) {
      addError('All fields are required')
      return
    }
    if (password !== confirmPassword) {
      addError('Passwords do not match')
      return
    }

    setIsSubmitting(true)
    try {
      await settingApi.settingAccount({ username: trimmedUsername, password })
      sessionStorage.setItem(PENDING_SUCCESS_MESSAGE_KEY, 'Account updated. Please sign in again.')
      localStorage.removeItem('token')
      navigate('/login', { replace: true })
    } catch (err) {
      addError(getErrorMessage(err, 'Failed to update account'))
      setIsSubmitting(false)
    }
  }

  return (
    <section className={panelStyles.tile}>
      <NotificationContainer errors={errors} successes={successes} onClose={removeNotification} />

      <div className={panelStyles.tileHeader}>
        <div>
          <p className={panelStyles.tileTag}>Account</p>
          <h3>Settings</h3>
          <p className={panelStyles.tileDescription}>
            Update your username and password. You'll need to sign in again afterward.
          </p>
        </div>
      </div>

      <form className={styles.form} onSubmit={handleSubmit}>
        <div className={modalStyles.formGroup}>
          <label className={modalStyles.label} htmlFor="setting-username">
            Username <span className={modalStyles.requiredMark}>*</span>
          </label>
          <input
            id="setting-username"
            className={modalStyles.input}
            value={username}
            onChange={(event) => setUsername(event.target.value)}
            autoComplete="username"
            disabled={isSubmitting}
            required
          />
        </div>

        <div className={modalStyles.formGroup}>
          <label className={modalStyles.label} htmlFor="setting-password">
            New Password <span className={modalStyles.requiredMark}>*</span>
          </label>
          <input
            id="setting-password"
            type="password"
            className={modalStyles.input}
            value={password}
            onChange={(event) => setPassword(event.target.value)}
            autoComplete="new-password"
            disabled={isSubmitting}
            required
          />
        </div>

        <div className={modalStyles.formGroup}>
          <label className={modalStyles.label} htmlFor="setting-confirm-password">
            Confirm Password <span className={modalStyles.requiredMark}>*</span>
          </label>
          <input
            id="setting-confirm-password"
            type="password"
            className={modalStyles.input}
            value={confirmPassword}
            onChange={(event) => setConfirmPassword(event.target.value)}
            autoComplete="new-password"
            disabled={isSubmitting}
            required
          />
        </div>

        <div className={styles.actionRow}>
          <Button type="submit" disabled={isSubmitting}>
            {isSubmitting ? 'Saving...' : 'Save Changes'}
          </Button>
        </div>
      </form>
    </section>
  )
}
