import { useState } from 'react'
import { useAuth } from '@/hooks/useAuth'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { PasswordInput } from '@/components/ui/password-input'
import { Label } from '@/components/ui/label'
import { PasswordRequirements } from '@/components/password-requirements'
import { validatePassword } from '@/lib/password-validation'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'

interface RegisterFormProps {
  onToggleToLogin: () => void
}

export function RegisterForm({ onToggleToLogin }: RegisterFormProps) {
  const { register, isLoading, error } = useAuth()
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [firstName, setFirstName] = useState('')
  const [lastName, setLastName] = useState('')
  const [vehicleName, setVehicleName] = useState('')
  const [vehicleMpg, setVehicleMpg] = useState('')

  async function handleSubmit(event: React.FormEvent) {
    event.preventDefault()

    const data = {
      username,
      password,
      first_name: firstName,
      last_name: lastName,
      ...(vehicleName ? { vehicle_name: vehicleName } : {}),
      ...(vehicleMpg ? { vehicle_mpg: parseFloat(vehicleMpg) } : {}),
    }

    await register(data)
  }

  return (
    <Card className="w-full max-w-md mx-auto">
      <CardHeader>
        <CardTitle>Create Account</CardTitle>
        <CardDescription>
          Register a new Grubbin Data account
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="register-username">Username</Label>
            <Input
              id="register-username"
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="Choose a username"
              required
              disabled={isLoading}
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="register-password">Password</Label>
            <PasswordInput
              id="register-password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Choose a password"
              required
              disabled={isLoading}
            />
            <PasswordRequirements password={password} />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="register-first-name">First Name</Label>
              <Input
                id="register-first-name"
                type="text"
                value={firstName}
                onChange={(e) => setFirstName(e.target.value)}
                placeholder="First name"
                required
                disabled={isLoading}
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="register-last-name">Last Name</Label>
              <Input
                id="register-last-name"
                type="text"
                value={lastName}
                onChange={(e) => setLastName(e.target.value)}
                placeholder="Last name"
                required
                disabled={isLoading}
              />
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="register-vehicle-name">Vehicle Name (optional)</Label>
            <Input
              id="register-vehicle-name"
              type="text"
              value={vehicleName}
              onChange={(e) => setVehicleName(e.target.value)}
              placeholder="e.g., Honda Civic"
              disabled={isLoading}
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="register-vehicle-mpg">Vehicle MPG (optional)</Label>
            <Input
              id="register-vehicle-mpg"
              type="number"
              step="0.1"
              value={vehicleMpg}
              onChange={(e) => setVehicleMpg(e.target.value)}
              placeholder="e.g., 32.5"
              disabled={isLoading}
            />
          </div>

          {error && (
            <p className="text-sm text-destructive">{error}</p>
          )}

          <Button
            type="submit"
            className="w-full"
            disabled={isLoading || !validatePassword(password).isValid}
          >
            {isLoading ? 'Creating account...' : 'Create Account'}
          </Button>

          <p className="text-sm text-center text-muted-foreground">
            Already have an account?{' '}
            <button
              type="button"
              onClick={onToggleToLogin}
              className="text-primary hover:underline"
            >
              Sign In
            </button>
          </p>
        </form>
      </CardContent>
    </Card>
  )
}
