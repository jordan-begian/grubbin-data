import { useAuth } from '@/hooks/useAuth'
import {
  Sidebar,
  SidebarContent,
  SidebarHeader,
  SidebarFooter,
  useSidebar,
} from '@/components/ui/sidebar'
import { Button } from '@/components/ui/button'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Separator } from '@/components/ui/separator'
import { User, Car, Calendar, LogOut } from 'lucide-react'

export function UserSidebar() {
  const { user, logout } = useAuth()
  const { open } = useSidebar()

  if (!user) return null

  const fullName = user.profile
    ? `${user.profile.first_name} ${user.profile.last_name}`
    : user.username

  const initials = user.profile
    ? `${user.profile.first_name[0]}${user.profile.last_name[0]}`
    : user.username.slice(0, 2).toUpperCase()

  const createdDate = new Date(user.created).toLocaleDateString()

  return (
    <Sidebar 
      side="right" 
      variant="sidebar"
      collapsible="icon"
      className="rounded-lg border bg-card shadow-sm"
    >
      <SidebarHeader className="p-4 border-b">
        <div className={`flex items-center ${open ? 'gap-3' : 'justify-center'}`}>
          <Avatar className="h-10 w-10 shrink-0">
            <AvatarFallback>{initials}</AvatarFallback>
          </Avatar>
          {open && (
            <div className="flex flex-col min-w-0 overflow-hidden">
              <span className="font-semibold truncate">{fullName}</span>
              <span className="text-xs text-muted-foreground truncate">
                {user.username}
              </span>
            </div>
          )}
        </div>
      </SidebarHeader>

      <SidebarContent className="p-4">
        {open && (
          <div className="space-y-6">
            <div className="space-y-3">
              <h4 className="text-sm font-medium text-muted-foreground">
                Profile
              </h4>
              <div className="space-y-2 text-sm">
                <div className="flex items-center gap-2">
                  <User className="h-4 w-4 text-muted-foreground shrink-0" />
                  <span className="truncate">
                    {user.profile?.first_name} {user.profile?.last_name}
                  </span>
                </div>
                <div className="flex items-center gap-2">
                  <Calendar className="h-4 w-4 text-muted-foreground shrink-0" />
                  <span>Member since {createdDate}</span>
                </div>
              </div>
            </div>

            <Separator />

            {user.profile?.vehicle && (
              <div className="space-y-3">
                <h4 className="text-sm font-medium text-muted-foreground">
                  Vehicle
                </h4>
                <div className="space-y-2 text-sm">
                  <div className="flex items-center gap-2">
                    <Car className="h-4 w-4 text-muted-foreground shrink-0" />
                    <span>{user.profile.vehicle.name}</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <span className="text-muted-foreground">MPG:</span>
                    <span>{user.profile.vehicle.average_mpg}</span>
                  </div>
                </div>
              </div>
            )}
          </div>
        )}
      </SidebarContent>

      <SidebarFooter className="p-4 border-t mt-auto">
        <Button 
          variant="outline" 
          className={`${open ? 'w-full' : 'w-full justify-center'}`}
          onClick={logout}
        >
          <LogOut className="h-4 w-4" />
          {open && <span className="ml-2">Sign Out</span>}
        </Button>
      </SidebarFooter>
    </Sidebar>
  )
}
