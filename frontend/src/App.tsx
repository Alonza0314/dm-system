import LoginPage from './page/login/LoginPage'
import { Navigate, Route, Routes } from 'react-router-dom'
import HomePage from './page/home/HomePage'
import CategoryDevicesPage from './page/categoryDevices/CategoryDevicesPage'
import DeviceDetailPage from './page/deviceDetail/DeviceDetailPage'
import QrCodePage from './page/qrcode/QrCodePage'
import AppLayout from './layout/AppLayout'

function RequireAuth({ children }: { children: React.ReactNode }) {
  const token = localStorage.getItem('token')
  if (!token) {
    return <Navigate to="/login" replace />
  }

  return <>{children}</>
}

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/qrcode/:cate/:dev" element={<QrCodePage />} />
      <Route
        element={(
          <RequireAuth>
            <AppLayout />
          </RequireAuth>
        )}
      >
        <Route path="/" element={<HomePage />} />
        <Route path="/category/:categoryName" element={<CategoryDevicesPage />} />
        <Route path="/category/:categoryName/device/:deviceName" element={<DeviceDetailPage />} />
      </Route>
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}
