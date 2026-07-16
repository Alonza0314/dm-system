import { useNavigate } from 'react-router-dom'
import Button from '../../components/button/button'
import { CategoryProvider, useCategoryContext } from '../../context/CategoryContext'
import CategoryPanel from './CategoryPanel'
import styles from './home-page.module.css'

function Sidebar({ onLogout }: { onLogout: () => void }) {
  const { categories } = useCategoryContext()

  return (
    <aside className={styles.sidebar}>
      <div>
        <p className={styles.badge}>DM System</p>
        <h1 className={styles.brand}>Dashboard</h1>

        <p className={styles.navLabel}>Menu</p>
        <nav className={styles.nav}>
          <a className={`${styles.navItem} ${styles.navItemActive}`} href="#">Category</a>
          {categories.map((category) => (
            <a key={category.name} className={styles.navItem} href="#">{category.name}</a>
          ))}
        </nav>
      </div>

      <div className={styles.logoutWrap}>
        <Button variant="secondary" onClick={onLogout}>Logout</Button>
      </div>
    </aside>
  )
}

export default function HomePage() {
  const navigate = useNavigate()

  function handleLogout() {
    localStorage.removeItem('token')
    navigate('/login', { replace: true })
  }

  return (
    <CategoryProvider>
      <div className={styles.layout}>
        <Sidebar onLogout={handleLogout} />

        <main className={styles.content}>
          <CategoryPanel />
        </main>
      </div>
    </CategoryProvider>
  )
}
