import { useCallback, useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import Button from '../../components/button/button'
import Modal from '../../components/modal/modal'
import NotificationContainer from '../../components/notifications/NotificationContainer'
import { useNotifications } from '../../hooks/useNotifications'
import { useCategoryContext } from '../../context/CategoryContext'
import { deviceApi } from '../../apiClient'
import { getErrorMessage } from '../../utils/getErrorMessage'
import { downloadDataUrl, generateQrCodeDataUrl } from '../../utils/qrCode'
import type { DeviceShort } from '../../api'
import styles from '../../styles/dashboard-panel.module.css'
import modalStyles from '../../components/modal/modal.module.css'

interface DeviceForm {
  category: string
  name: string
  owner: string
  note: string
}

interface QrModalState {
  deviceName: string
  url: string
  dataUrl: string
}

export default function CategoryDevicesPage() {
  const { categoryName = '' } = useParams()
  const navigate = useNavigate()
  const { categories } = useCategoryContext()
  const { errors, successes, addError, addSuccess, removeNotification } = useNotifications()

  const [devices, setDevices] = useState<DeviceShort[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const [isCreateOpen, setCreateOpen] = useState(false)
  const [form, setForm] = useState<DeviceForm>({ category: categoryName, name: '', owner: '', note: '' })
  const [isSubmitting, setSubmitting] = useState(false)

  const [confirmTarget, setConfirmTarget] = useState<string | null>(null)
  const [pendingDelete, setPendingDelete] = useState<string | null>(null)

  const [qrModal, setQrModal] = useState<QrModalState | null>(null)

  const fetchDevices = useCallback(async (silent = false) => {
    if (!silent) setIsLoading(true)
    try {
      const response = await deviceApi.listDevices(categoryName)
      setDevices(response.data.Devices ?? [])
      setError(null)
    } catch (err) {
      setError(getErrorMessage(err, 'Failed to load devices'))
    } finally {
      if (!silent) setIsLoading(false)
    }
  }, [categoryName])

  useEffect(() => {
    fetchDevices()
    const interval = setInterval(() => fetchDevices(true), 10_000)
    return () => clearInterval(interval)
  }, [fetchDevices])

  function openCreate() {
    setForm({ category: categoryName, name: '', owner: '', note: '' })
    setCreateOpen(true)
  }

  async function handleCreate() {
    const name = form.name.trim()
    if (!form.category || !name) {
      addError('Category and name are required')
      return
    }

    setSubmitting(true)
    try {
      await deviceApi.createDevice({
        category: form.category,
        name,
        owner: form.owner.trim() || undefined,
        note: form.note.trim() || undefined,
      })
      addSuccess(`Device "${name}" created`)
      setCreateOpen(false)
      await fetchDevices()
    } catch (err) {
      addError(getErrorMessage(err, 'Failed to create device'))
    } finally {
      setSubmitting(false)
    }
  }

  async function handleConfirmDelete() {
    if (!confirmTarget) return

    setPendingDelete(confirmTarget)
    try {
      await deviceApi.deleteDevice(categoryName, confirmTarget)
      addSuccess(`Device "${confirmTarget}" deleted`)
      setConfirmTarget(null)
      await fetchDevices()
    } catch (err) {
      addError(getErrorMessage(err, 'Failed to delete device'))
    } finally {
      setPendingDelete(null)
    }
  }

  function goToDetail(deviceName: string) {
    navigate(`/category/${encodeURIComponent(categoryName)}/device/${encodeURIComponent(deviceName)}`)
  }

  async function handleShowQr(deviceName: string) {
    const url = `${window.location.origin}/qrcode/${encodeURIComponent(categoryName)}/${encodeURIComponent(deviceName)}`
    try {
      const dataUrl = await generateQrCodeDataUrl(url)
      setQrModal({ deviceName, url, dataUrl })
    } catch (err) {
      addError(getErrorMessage(err, 'Failed to generate QR code'))
    }
  }

  function handleDownloadQr() {
    if (!qrModal) return
    downloadDataUrl(qrModal.dataUrl, `${categoryName}-${qrModal.deviceName}-qrcode.png`)
  }

  return (
    <section className={styles.tile}>
      <NotificationContainer errors={errors} successes={successes} onClose={removeNotification} />

      <div className={styles.tileHeader}>
        <div>
          <p className={styles.tileTag}>Category</p>
          <h3>{categoryName}</h3>
          <p className={styles.tileDescription}>Devices registered under this category.</p>
        </div>
        <Button onClick={openCreate}>Add Device</Button>
      </div>

      {error && <p className={styles.tableError}>{error}</p>}

      <div className={styles.tableWrap}>
        <table className={styles.table}>
          <thead>
            <tr>
              <th>Name</th>
              <th>Status</th>
              <th>User</th>
              <th aria-label="Actions" />
            </tr>
          </thead>
          <tbody>
            {devices.map((device) => (
              <tr key={device.name}>
                <td>{device.name}</td>
                <td>
                  <span
                    className={`${styles.statusPill} ${
                      device.status === 'using' ? styles.statusUsing : styles.statusIdle
                    }`}
                  >
                    {device.status}
                  </span>
                </td>
                <td>{device.user || '—'}</td>
                <td className={styles.tableActions}>
                  <Button
                    variant="utility"
                    className={styles.qrButtonSpacing}
                    onClick={() => handleShowQr(device.name ?? '')}
                  >
                    QR Code
                  </Button>
                  <Button variant="secondary" onClick={() => goToDetail(device.name ?? '')}>
                    Detail
                  </Button>
                  <Button
                    variant="secondary"
                    onClick={() => setConfirmTarget(device.name ?? '')}
                    disabled={pendingDelete === device.name}
                  >
                    {pendingDelete === device.name ? 'Deleting...' : 'Delete'}
                  </Button>
                </td>
              </tr>
            ))}
            {!isLoading && devices.length === 0 && (
              <tr>
                <td colSpan={4} className={styles.tableEmpty}>No devices yet.</td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      <Modal
        isOpen={isCreateOpen}
        onClose={() => setCreateOpen(false)}
        title="Add Device"
        onSubmit={handleCreate}
        submitLabel={isSubmitting ? 'Creating...' : 'Create'}
        isSubmitDisabled={isSubmitting}
      >
        <div className={modalStyles.formGroup}>
          <label className={modalStyles.label} htmlFor="device-category">Category</label>
          <select
            id="device-category"
            className={modalStyles.input}
            value={form.category}
            onChange={(event) => setForm((prev) => ({ ...prev, category: event.target.value }))}
            disabled={isSubmitting}
          >
            {categories.map((category) => (
              <option key={category.name} value={category.name}>{category.name}</option>
            ))}
          </select>
        </div>

        <div className={modalStyles.formGroup}>
          <label className={modalStyles.label} htmlFor="device-name">
            Name <span className={modalStyles.requiredMark}>*</span>
          </label>
          <input
            id="device-name"
            className={modalStyles.input}
            value={form.name}
            onChange={(event) => setForm((prev) => ({ ...prev, name: event.target.value }))}
            disabled={isSubmitting}
            autoFocus
          />
        </div>

        <div className={modalStyles.formGroup}>
          <label className={modalStyles.label} htmlFor="device-owner">
            Owner <span className={modalStyles.optionalHint}>(optional)</span>
          </label>
          <input
            id="device-owner"
            className={modalStyles.input}
            value={form.owner}
            onChange={(event) => setForm((prev) => ({ ...prev, owner: event.target.value }))}
            disabled={isSubmitting}
          />
        </div>

        <div className={modalStyles.formGroup}>
          <label className={modalStyles.label} htmlFor="device-note">
            Note <span className={modalStyles.optionalHint}>(optional)</span>
          </label>
          <input
            id="device-note"
            className={modalStyles.input}
            value={form.note}
            onChange={(event) => setForm((prev) => ({ ...prev, note: event.target.value }))}
            disabled={isSubmitting}
          />
        </div>
      </Modal>

      <Modal
        isOpen={confirmTarget !== null}
        onClose={() => setConfirmTarget(null)}
        title="Delete Device"
        onSubmit={handleConfirmDelete}
        submitLabel={pendingDelete ? 'Deleting...' : 'Delete'}
        isSubmitDisabled={pendingDelete !== null}
      >
        <p>
          Are you sure you want to delete device "{confirmTarget}"? This action cannot be undone.
        </p>
      </Modal>

      <Modal
        isOpen={qrModal !== null}
        onClose={() => setQrModal(null)}
        title="Device QR Code"
        onSubmit={handleDownloadQr}
        submitLabel="Download"
      >
        {qrModal && (
          <div className={styles.qrPreview}>
            <img
              src={qrModal.dataUrl}
              alt={`QR code for ${qrModal.deviceName}`}
              className={styles.qrImage}
            />
            <a href={qrModal.url} target="_blank" rel="noopener noreferrer" className={styles.qrLink}>
              {qrModal.url}
            </a>
          </div>
        )}
      </Modal>
    </section>
  )
}
