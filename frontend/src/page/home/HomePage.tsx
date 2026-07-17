import { useCategoryContext } from '../../context/CategoryContext'
import StatusOverview from './StatusOverview'
import CategoryPanel from './CategoryPanel'

export default function HomePage() {
  const { categories } = useCategoryContext()

  return (
    <>
      <StatusOverview categories={categories} />
      <CategoryPanel />
    </>
  )
}
