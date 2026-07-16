import { useNavigate } from 'react-router-dom'
import Button from '../../components/button/button'
import styles from './home-page.module.css'

export default function HomePage() {
  const navigate = useNavigate()

  function handleLogout() {
    localStorage.removeItem('token')
    navigate('/login', { replace: true })
  }

  return (
    <div className={styles.layout}>
      <aside className={styles.sidebar}>
        <div>
          <p className={styles.badge}>DM System</p>
          <h1 className={styles.brand}>Console</h1>

          <p className={styles.navLabel}>Menu</p>
          <nav className={styles.nav}>
            <a className={`${styles.navItem} ${styles.navItemActive}`} href="#">Home</a>
            <a className={styles.navItem} href="#">Module A</a>
            <a className={styles.navItem} href="#">Module B</a>
          </nav>
        </div>

        <div className={styles.logoutWrap}>
          <Button variant="secondary" onClick={handleLogout}>Logout</Button>
        </div>
      </aside>

      <main className={styles.content}>
        <section className={styles.heroTile}>
          <p className={styles.eyebrow}>Console</p>
          <h2 className={styles.heroTitle}>Home</h2>
          <p className={styles.heroSubtitle}>A clean framework canvas ready for your features.</p>
        </section>

        <section className={styles.tile}>
          <div className={styles.tileInner}>
            <div className={styles.tileText}>
              <p className={styles.tileTag}>Widget Area</p>
              <h3>Place dashboard cards</h3>
              <p>Drop in tables, charts, or summary cards here as the first module goes in.</p>
            </div>
            <div className={styles.tileGraphic} aria-hidden="true" />
          </div>
        </section>

        <section className={`${styles.tile} ${styles.tileParchment}`}>
          <div className={styles.tileInner}>
            <div className={styles.tileText}>
              <p className={styles.tileTag}>Feature Area</p>
              <h3>Start your first module</h3>
              <p>Use this section as a starting point for the next module page.</p>
            </div>
            <div className={styles.tileGraphic} aria-hidden="true" />
          </div>
        </section>
      </main>
    </div>
  )
}
