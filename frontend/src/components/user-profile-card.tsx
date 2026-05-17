import { useAuth } from '@/hooks/useAuth'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'

export function UserProfileCard() {
  const { user, logout } = useAuth()

  if (!user) return null

  const fullName = user.profile
    ? `${user.profile.first_name} ${user.profile.last_name}`
    : user.username

  const createdDate = new Date(user.created).toLocaleDateString()

  return (
    <Card className="w-full max-w-md mx-auto">
      <CardHeader>
        <CardTitle>Welcome, {fullName}</CardTitle>
        <CardDescription>
          Your Grubbin Data profile
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="space-y-2">
          <div className="flex justify-between">
            <span className="text-muted-foreground">Username</span>
            <span className="font-medium">{user.username}</span>
          </div>

          {user.profile && (
            <>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Name</span>
                <span className="font-medium">
                  {user.profile.first_name} {user.profile.last_name}
                </span>
              </div>

              {user.profile.vehicle && (
                <>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Vehicle</span>
                    <span className="font-medium">{user.profile.vehicle.name}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">MPG</span>
                    <span className="font-medium">{user.profile.vehicle.average_mpg}</span>
                  </div>
                </>
              )}
            </>
          )}

          <div className="flex justify-between">
            <span className="text-muted-foreground">Member Since</span>
            <span className="font-medium">{createdDate}</span>
          </div>
        </div>

        <Button variant="outline" className="w-full" onClick={logout}>
          Sign Out
        </Button>
      </CardContent>
    </Card>
  )
}
