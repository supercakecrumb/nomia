import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import './index.css'
import './i18n/config'
import AppLayout from './components/layout/AppLayout'
import MainPage from './pages/MainPage'
import NamesExplorerPage from './pages/NamesExplorerPage'
import NameDetailPage from './pages/NameDetailPage'

// Configure React Query with sensible defaults
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 5, // 5 minutes - data is considered fresh for 5 minutes
      gcTime: 1000 * 60 * 30, // 30 minutes - cache garbage collection time (was cacheTime)
      retry: 2, // Retry failed requests twice
      refetchOnWindowFocus: false, // Don't refetch on window focus by default
    },
  },
})

// Configure routes with React Router 7
const router = createBrowserRouter([
  {
    path: '/',
    element: <AppLayout />,
    children: [
      {
        index: true,
        element: <MainPage />,
      },
      {
        path: 'names',
        element: <NamesExplorerPage />,
      },
      {
        path: 'name/:name',
        element: <NameDetailPage />,
      },
    ],
  },
  {
    path: '*',
    element: (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <h1 className="text-6xl font-bold text-gray-900 mb-4">404</h1>
          <p className="text-xl text-gray-600 mb-8">Page not found</p>
          <a href="/" className="text-blue-600 hover:text-blue-800 font-medium">
            Go back home
          </a>
        </div>
      </div>
    ),
  },
])

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>
  </StrictMode>,
)
