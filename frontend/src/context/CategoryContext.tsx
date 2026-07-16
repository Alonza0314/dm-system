import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from 'react'
import type { Category } from '../api'
import { categoryApi } from '../apiClient'
import { getErrorMessage } from '../utils/getErrorMessage'

const POLL_INTERVAL_MS = 10_000

interface CategoryContextValue {
  categories: Category[]
  isLoading: boolean
  error: string | null
  createCategory: (name: string) => Promise<void>
  deleteCategory: (name: string) => Promise<void>
}

const CategoryContext = createContext<CategoryContextValue | null>(null)

function sortByNameDesc(categories: Category[]): Category[] {
  return [...categories].sort((a, b) => (b.name ?? '').localeCompare(a.name ?? ''))
}

export function CategoryProvider({ children }: { children: ReactNode }) {
  const [categories, setCategories] = useState<Category[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchCategories = useCallback(async (silent = false) => {
    if (!silent) setIsLoading(true)
    try {
      const response = await categoryApi.listCategories()
      setCategories(sortByNameDesc(response.data.categories ?? []))
      setError(null)
    } catch (err) {
      setError(getErrorMessage(err, 'Failed to load categories'))
    } finally {
      if (!silent) setIsLoading(false)
    }
  }, [])

  useEffect(() => {
    fetchCategories()
    const interval = setInterval(() => fetchCategories(true), POLL_INTERVAL_MS)
    return () => clearInterval(interval)
  }, [fetchCategories])

  const createCategory = useCallback(async (name: string) => {
    await categoryApi.createCategory({ name })
    await fetchCategories(true)
  }, [fetchCategories])

  const deleteCategory = useCallback(async (name: string) => {
    await categoryApi.deleteCategory(name)
    await fetchCategories(true)
  }, [fetchCategories])

  const value = useMemo(
    () => ({ categories, isLoading, error, createCategory, deleteCategory }),
    [categories, isLoading, error, createCategory, deleteCategory],
  )

  return <CategoryContext.Provider value={value}>{children}</CategoryContext.Provider>
}

export function useCategoryContext() {
  const context = useContext(CategoryContext)
  if (!context) {
    throw new Error('useCategoryContext must be used within a CategoryProvider')
  }
  return context
}
