import { useState } from 'react'
import { useAuth } from '@/hooks/useAuth'
import {
  SidebarProvider,
  SidebarTrigger,
} from '@/components/ui/sidebar'
import { ThemeSwitcher } from '@/components/ThemeSwitcher'
import { LoginForm } from '@/components/login-form'
import { RegisterForm } from '@/components/register-form'
import { UserSidebar } from '@/components/user-sidebar'
import { DeliveryForm } from '@/components/delivery-form'
import { DeliveryList } from '@/components/delivery-list'
import { DeliveryStatsCard } from '@/components/delivery-stats'
import { useDeliveries } from '@/hooks/useDeliveries'
import type {
  CreateDeliveryRequest,
  UpdateDeliveryRequest,
} from '@/types/delivery'
import { Menu } from 'lucide-react'

function AuthenticatedDashboard() {
  const { user } = useAuth()
  const {
    deliveries,
    stats,
    isLoading,
    createDelivery,
    updateDeliveries,
    deleteDeliveries,
    isCreating,
    isUpdating,
    isDeleting,
  } = useDeliveries(user?.id)

  const handleCreateDelivery = async (data: CreateDeliveryRequest) => {
    await createDelivery(data)
  }

  const handleEditDelivery = async (updated: UpdateDeliveryRequest) => {
    await updateDeliveries([updated])
  }

  const handleDeleteDeliveries = async (ids: string[]) => {
    await deleteDeliveries(ids)
  }

  return (
    <div className="flex min-h-svh w-full flex-col">
      {/* Header - full width */}
      <div className="flex items-center justify-between p-4 border-b shrink-0">
        <h1 className="text-xl font-semibold">Grubbin Data Dashboard</h1>
        <div className="flex items-center gap-4">
          <ThemeSwitcher />
          <SidebarTrigger>
            <Menu className="h-5 w-5" />
          </SidebarTrigger>
        </div>
      </div>

      {/* Content area - sidebar + main content side by side */}
      <div className="flex flex-1 overflow-hidden p-6 gap-6">
        {/* Main content */}
        <div className="flex-1 overflow-auto space-y-6">
          <DeliveryStatsCard
            stats={stats}
            deliveryCount={deliveries.length}
            isLoading={isLoading}
          />

          <div className="grid grid-cols-1 xl:grid-cols-2 gap-6">
            <DeliveryForm
              onSubmit={handleCreateDelivery}
              isSubmitting={isCreating}
            />

            <DeliveryList
              deliveries={deliveries}
              isLoading={isLoading}
              onEdit={handleEditDelivery}
              onDelete={handleDeleteDeliveries}
              isEditing={isUpdating}
              isDeleting={isDeleting}
            />
          </div>
        </div>

        {/* Sidebar - in same container, naturally aligned */}
        <UserSidebar />
      </div>
    </div>
  )
}

export default function App() {
  const { isAuthenticated } = useAuth()
  const [showRegister, setShowRegister] = useState(false)

  if (!isAuthenticated) {
    return (
      <main className="min-h-svh flex flex-col items-center justify-center p-4">
        <ThemeSwitcher />
        {showRegister ? (
          <RegisterForm onToggleToLogin={() => setShowRegister(false)} />
        ) : (
          <LoginForm onToggleToRegister={() => setShowRegister(true)} />
        )}
      </main>
    )
  }

  return (
    <SidebarProvider defaultOpen={false}>
      <AuthenticatedDashboard />
    </SidebarProvider>
  )
}
