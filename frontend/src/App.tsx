import { useState } from 'react'
import { useAuth } from '@/hooks/useAuth'
import { ThemeSwitcher } from '@/components/ThemeSwitcher'
import { LoginForm } from '@/components/login-form'
import { RegisterForm } from '@/components/register-form'
import { UserProfileCard } from '@/components/user-profile-card'

export default function App() {
  const { isAuthenticated } = useAuth()
  const [showRegister, setShowRegister] = useState(false)

  return (
    <main className="min-h-svh flex flex-col">
      <ThemeSwitcher />
      <div className="flex-1 flex items-center justify-center p-4">
        {isAuthenticated ? (
          <UserProfileCard />
        ) : showRegister ? (
          <RegisterForm onToggleToLogin={() => setShowRegister(false)} />
        ) : (
          <LoginForm onToggleToRegister={() => setShowRegister(true)} />
        )}
      </div>
    </main>
  )
}
