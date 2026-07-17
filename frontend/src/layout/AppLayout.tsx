import { NavLink, Outlet, useNavigate } from 'react-router-dom'
import Button from '../components/button/button'
import { CategoryProvider, useCategoryContext } from '../context/CategoryContext'
import styles from './app-layout.module.css'

function navItemClassName({ isActive }: { isActive: boolean }) {
  return `${styles.navItem} ${isActive ? styles.navItemActive : ''}`
}

function navSubItemClassName({ isActive }: { isActive: boolean }) {
  return `${styles.navSubItem} ${isActive ? styles.navSubItemActive : ''}`
}

function Sidebar({ onLogout }: { onLogout: () => void }) {
  const { categories } = useCategoryContext()

  return (
    <aside className={styles.sidebar}>
      <div>
        <p className={styles.badge}>DM System</p>
        <h1 className={styles.brand}>Console</h1>

        <p className={styles.navLabel}>Menu</p>
        <nav className={styles.nav}>
          <NavLink to="/" end className={navItemClassName}>Category</NavLink>
        </nav>

        {categories.length > 0 && (
          <>
            <p className={styles.navSubLabel}>Categories</p>
            <nav className={styles.nav}>
              {categories.map((category) => (
                <NavLink
                  key={category.name}
                  to={`/category/${encodeURIComponent(category.name ?? '')}`}
                  className={navSubItemClassName}
                >
                  {category.name}
                </NavLink>
              ))}
            </nav>
          </>
        )}
      </div>

      <div className={styles.logoutWrap}>
        <Button variant="secondary" onClick={onLogout}>Logout</Button>
      </div>
    </aside>
  )
}

export default function AppLayout() {
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
          <Outlet />
        </main>
      </div>
    </CategoryProvider>
  )
}
