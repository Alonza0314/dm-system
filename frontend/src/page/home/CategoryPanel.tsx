import { useState } from 'react'
import Button from '../../components/button/button'
import Modal from '../../components/modal/modal'
import NotificationContainer from '../../components/notifications/NotificationContainer'
import { useNotifications } from '../../hooks/useNotifications'
import { useCategoryContext } from '../../context/CategoryContext'
import { getErrorMessage } from '../../utils/getErrorMessage'
import styles from './home-page.module.css'
import modalStyles from '../../components/modal/modal.module.css'

export default function CategoryPanel() {
  const { categories, isLoading, error, createCategory, deleteCategory } = useCategoryContext()
  const { errors, successes, addError, addSuccess, removeNotification } = useNotifications()

  const [isCreateOpen, setCreateOpen] = useState(false)
  const [name, setName] = useState('')
  const [isSubmitting, setSubmitting] = useState(false)
  const [confirmTarget, setConfirmTarget] = useState<string | null>(null)
  const [pendingDelete, setPendingDelete] = useState<string | null>(null)

  function openCreate() {
    setName('')
    setCreateOpen(true)
  }

  async function handleCreate() {
    const trimmed = name.trim()
    if (!trimmed) {
      addError('Name is required')
      return
    }

    setSubmitting(true)
    try {
      await createCategory(trimmed)
      addSuccess(`Category "${trimmed}" created`)
      setCreateOpen(false)
    } catch (err) {
      addError(getErrorMessage(err, 'Failed to create category'))
    } finally {
      setSubmitting(false)
    }
  }

  async function handleConfirmDelete() {
    if (!confirmTarget) return

    setPendingDelete(confirmTarget)
    try {
      await deleteCategory(confirmTarget)
      addSuccess(`Category "${confirmTarget}" deleted`)
      setConfirmTarget(null)
    } catch (err) {
      addError(getErrorMessage(err, 'Failed to delete category'))
    } finally {
      setPendingDelete(null)
    }
  }

  return (
    <section className={styles.tile}>
      <NotificationContainer errors={errors} successes={successes} onClose={removeNotification} />

      <div className={styles.tileHeader}>
        <div>
          <p className={styles.tileTag}>Categories</p>
          <h3>Device Categories</h3>
        </div>
        <Button onClick={openCreate}>Create Category</Button>
      </div>

      {error && <p className={styles.tableError}>{error}</p>}

      <div className={styles.tableWrap}>
        <table className={styles.table}>
          <thead>
            <tr>
              <th>Name</th>
              <th>Idle</th>
              <th>Using</th>
              <th aria-label="Actions" />
            </tr>
          </thead>
          <tbody>
            {categories.map((category) => (
              <tr key={category.name}>
                <td>{category.name}</td>
                <td>{category.idle_device}</td>
                <td>{category.using_device}</td>
                <td className={styles.tableActions}>
                  <Button
                    variant="secondary"
                    onClick={() => setConfirmTarget(category.name ?? '')}
                    disabled={pendingDelete === category.name}
                  >
                    {pendingDelete === category.name ? 'Deleting...' : 'Delete'}
                  </Button>
                </td>
              </tr>
            ))}
            {!isLoading && categories.length === 0 && (
              <tr>
                <td colSpan={4} className={styles.tableEmpty}>No categories yet.</td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      <Modal
        isOpen={isCreateOpen}
        onClose={() => setCreateOpen(false)}
        title="Create Category"
        onSubmit={handleCreate}
      >
        <div className={modalStyles.formGroup}>
          <label className={modalStyles.label} htmlFor="category-name">Name</label>
          <input
            id="category-name"
            className={modalStyles.input}
            value={name}
            onChange={(event) => setName(event.target.value)}
            disabled={isSubmitting}
            autoFocus
          />
        </div>
      </Modal>

      <Modal
        isOpen={confirmTarget !== null}
        onClose={() => setConfirmTarget(null)}
        title="Delete Category"
        onSubmit={handleConfirmDelete}
        submitLabel={pendingDelete ? 'Deleting...' : 'Delete'}
        isSubmitDisabled={pendingDelete !== null}
      >
        <p>
          Are you sure you want to delete category "{confirmTarget}"? This action cannot be undone.
        </p>
      </Modal>
    </section>
  )
}
