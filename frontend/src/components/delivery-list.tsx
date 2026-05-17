import { useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Skeleton } from '@/components/ui/skeleton'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { DeliveryEditForm } from './delivery-edit-form'
import type { DeliveryResponse, UpdateDeliveryRequest } from '@/types/delivery'
import {
  MoreVertical,
  Trash2,
  Edit,
  MapPin,
  DollarSign,
  Clock,
  Calendar,
  Loader2,
} from 'lucide-react'

interface DeliveryListProps {
  deliveries: DeliveryResponse[]
  isLoading: boolean
  onEdit: (delivery: UpdateDeliveryRequest) => Promise<void>
  onDelete: (deliveryIds: string[]) => Promise<void>
  isEditing: boolean
  isDeleting: boolean
}

export function DeliveryList({
  deliveries,
  isLoading,
  onEdit,
  onDelete,
  isEditing,
  isDeleting,
}: DeliveryListProps) {
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set())
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [singleDeleteId, setSingleDeleteId] = useState<string | null>(null)

  const toggleSelection = (id: string) => {
    const newSelected = new Set(selectedIds)
    if (newSelected.has(id)) {
      newSelected.delete(id)
    } else {
      newSelected.add(id)
    }
    setSelectedIds(newSelected)
  }

  const toggleAll = () => {
    if (selectedIds.size === deliveries.length) {
      setSelectedIds(new Set())
    } else {
      setSelectedIds(new Set(deliveries.map((d) => d.id)))
    }
  }

  const handleBulkDelete = async () => {
    await onDelete(Array.from(selectedIds))
    setSelectedIds(new Set())
    setDeleteDialogOpen(false)
  }

  const handleSingleDelete = async () => {
    if (singleDeleteId) {
      await onDelete([singleDeleteId])
      setSingleDeleteId(null)
    }
  }

  const formatCurrency = (cents: number) => {
    return `$${(cents / 100).toFixed(2)}`
  }

  const formatDuration = (start: string, end: string) => {
    const duration = new Date(end).getTime() - new Date(start).getTime()
    const minutes = Math.floor(duration / 60000)
    return `${minutes} min`
  }

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <Skeleton className="h-8 w-48" />
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {[1, 2, 3].map((i) => (
              <div key={i} className="border rounded-lg p-4 space-y-3">
                <div className="flex items-center gap-4">
                  <Skeleton className="h-4 w-4" />
                  <div className="flex-1 space-y-2">
                    <Skeleton className="h-4 w-full" />
                    <Skeleton className="h-4 w-3/4" />
                  </div>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    )
  }

  if (deliveries.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Your Deliveries</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-center py-12 text-muted-foreground">
            <div className="mb-4">
              <MapPin className="h-12 w-12 mx-auto text-muted-foreground/50" />
            </div>
            <p className="text-lg font-medium">No deliveries yet</p>
            <p className="text-sm">
              Create your first delivery using the form!
            </p>
          </div>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <CardTitle>Your Deliveries ({deliveries.length})</CardTitle>
        {selectedIds.size > 0 && (
          <Button
            variant="destructive"
            size="sm"
            onClick={() => setDeleteDialogOpen(true)}
            disabled={isDeleting}
          >
            {isDeleting ? (
              <Loader2 className="h-4 w-4 mr-2 animate-spin" />
            ) : (
              <Trash2 className="h-4 w-4 mr-2" />
            )}
            Delete ({selectedIds.size})
          </Button>
        )}
      </CardHeader>
      <CardContent>
        <ScrollArea className="h-[600px]">
          <div className="space-y-4">
            <div className="flex items-center gap-4 px-2 py-2 bg-muted rounded-lg">
              <Checkbox
                checked={
                  selectedIds.size === deliveries.length &&
                  deliveries.length > 0
                }
                onCheckedChange={toggleAll}
              />
              <span className="text-sm font-medium flex-1">Select All</span>
            </div>

            {deliveries.map((delivery) => (
              <div
                key={delivery.id}
                className="border rounded-lg p-4 space-y-3 hover:bg-muted/50 transition-colors"
              >
                <div className="flex items-start gap-4">
                  <Checkbox
                    checked={selectedIds.has(delivery.id)}
                    onCheckedChange={() => toggleSelection(delivery.id)}
                  />

                  <div className="flex-1 space-y-2">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        <Calendar className="h-4 w-4 text-muted-foreground" />
                        <span className="text-sm text-muted-foreground">
                          {new Date(delivery.created).toLocaleDateString()}
                        </span>
                        <Clock className="h-4 w-4 text-muted-foreground ml-2" />
                        <span className="text-sm text-muted-foreground">
                          {formatDuration(delivery.start, delivery.end)}
                        </span>
                      </div>

                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="sm">
                            <MoreVertical className="h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <Dialog>
                            <DialogTrigger asChild>
                              <DropdownMenuItem
                                onSelect={(e: Event) => e.preventDefault()}
                              >
                                <Edit className="h-4 w-4 mr-2" />
                                Edit
                              </DropdownMenuItem>
                            </DialogTrigger>
                            <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
                              <DialogHeader>
                                <DialogTitle>Edit Delivery</DialogTitle>
                              </DialogHeader>
                              <DeliveryEditForm
                                delivery={delivery}
                                onSubmit={async (updated) => {
                                  await onEdit(updated)
                                }}
                                isSubmitting={isEditing}
                              />
                            </DialogContent>
                          </Dialog>

                          <DropdownMenuItem
                            className="text-destructive"
                            onSelect={(e: Event) => {
                              e.preventDefault()
                              setSingleDeleteId(delivery.id)
                            }}
                          >
                            <Trash2 className="h-4 w-4 mr-2" />
                            Delete
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </div>

                    <div className="flex items-center gap-2 text-sm">
                      <MapPin className="h-4 w-4 text-primary" />
                      <span className="font-medium">
                        {delivery.pickup.name}
                      </span>
                      <span className="text-muted-foreground">→</span>
                      <span>
                        Dropoff ({delivery.dropoff.lat.toFixed(4)},{' '}
                        {delivery.dropoff.lon.toFixed(4)})
                      </span>
                    </div>

                    <div className="flex items-center gap-4 flex-wrap">
                      <div className="flex items-center gap-1">
                        <DollarSign className="h-4 w-4 text-green-600" />
                        <span className="font-semibold">
                          {formatCurrency(
                            delivery.earnings.base +
                              delivery.earnings.tip +
                              (delivery.earnings.bonus || 0)
                          )}
                        </span>
                      </div>

                      <div className="flex gap-2 flex-wrap">
                        <Badge variant="secondary">
                          Base: {formatCurrency(delivery.earnings.base)}
                        </Badge>
                        <Badge variant="secondary">
                          Tip: {formatCurrency(delivery.earnings.tip)}
                        </Badge>
                        {delivery.earnings.bonus && (
                          <Badge variant="secondary">
                            Bonus:{' '}
                            {formatCurrency(delivery.earnings.bonus)}
                          </Badge>
                        )}
                      </div>
                    </div>

                    {delivery.note && (
                      <p className="text-sm text-muted-foreground italic border-l-2 border-muted pl-3">
                        &ldquo;{delivery.note}&rdquo;
                      </p>
                    )}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </ScrollArea>
      </CardContent>

      <AlertDialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Deliveries</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete {selectedIds.size} delivery
              {selectedIds.size > 1 ? 'ies' : 'y'}? This action cannot be
              undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel onClick={() => setDeleteDialogOpen(false)}>
              Cancel
            </AlertDialogCancel>
            <AlertDialogAction
              onClick={handleBulkDelete}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
              disabled={isDeleting}
            >
              {isDeleting ? (
                <>
                  <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                  Deleting...
                </>
              ) : (
                'Delete'
              )}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <AlertDialog
        open={!!singleDeleteId}
        onOpenChange={() => setSingleDeleteId(null)}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Delivery</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete this delivery? This action cannot
              be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel onClick={() => setSingleDeleteId(null)}>
              Cancel
            </AlertDialogCancel>
            <AlertDialogAction
              onClick={handleSingleDelete}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
              disabled={isDeleting}
            >
              {isDeleting ? (
                <>
                  <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                  Deleting...
                </>
              ) : (
                'Delete'
              )}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </Card>
  )
}
