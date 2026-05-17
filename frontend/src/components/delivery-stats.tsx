import { Card, CardContent } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import type { DeliveryStats } from '@/types/delivery'
import {
  Package,
  DollarSign,
  Clock,
  Navigation,
  Fuel,
  TrendingUp,
} from 'lucide-react'

interface DeliveryStatsProps {
  stats?: DeliveryStats
  deliveryCount: number
  isLoading: boolean
}

export function DeliveryStatsCard({
  stats,
  deliveryCount,
  isLoading,
}: DeliveryStatsProps) {
  if (isLoading) {
    return (
      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
        {[1, 2, 3, 4, 5, 6].map((i) => (
          <Card key={i}>
            <CardContent className="p-4">
              <Skeleton className="h-8 w-24 mb-2" />
              <Skeleton className="h-4 w-16" />
            </CardContent>
          </Card>
        ))}
      </div>
    )
  }

  if (!stats || deliveryCount === 0) {
    return null
  }

  const totalEarnings =
    ((stats.average_tip + stats.average_base_pay) * deliveryCount) / 100

  const statItems = [
    {
      icon: Package,
      label: 'Deliveries',
      value: deliveryCount.toString(),
    },
    {
      icon: DollarSign,
      label: 'Total Earnings',
      value: `$${totalEarnings.toFixed(2)}`,
    },
    {
      icon: TrendingUp,
      label: 'Avg per Delivery',
      value: `$${(
        (stats.average_tip + stats.average_base_pay) /
        100
      ).toFixed(2)}`,
    },
    {
      icon: Clock,
      label: 'Avg Time',
      value: `${Math.round(stats.average_delivery_time / 60)} min`,
    },
    {
      icon: Navigation,
      label: 'Total Miles',
      value: stats.total_miles.toFixed(1),
    },
    {
      icon: Fuel,
      label: 'Est. Fuel Cost',
      value: stats.used_fuel_cost
        ? `$${(stats.used_fuel_cost / 100).toFixed(2)}`
        : 'N/A',
    },
  ]

  return (
    <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
      {statItems.map((item) => (
        <Card key={item.label}>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 text-muted-foreground mb-2">
              <item.icon className="h-4 w-4" />
              <span className="text-xs">{item.label}</span>
            </div>
            <div className="text-2xl font-bold">{item.value}</div>
          </CardContent>
        </Card>
      ))}
    </div>
  )
}
