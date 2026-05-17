import { useState } from 'react'
import { format } from 'date-fns'
import { Calendar as CalendarIcon } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Separator } from '@/components/ui/separator'
import { Label } from '@/components/ui/label'
import { Calendar } from '@/components/ui/calendar'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'
import { Loader2 } from 'lucide-react'
import type { CreateDeliveryRequest } from '@/types/delivery'

interface DeliveryFormProps {
  onSubmit: (data: CreateDeliveryRequest) => Promise<void>
  isSubmitting: boolean
}

export function DeliveryForm({ onSubmit, isSubmitting }: DeliveryFormProps) {
  const [startDate, setStartDate] = useState<Date>(new Date())
  const [endDate, setEndDate] = useState<Date>(new Date())
  const [startTime, setStartTime] = useState('12:00')
  const [endTime, setEndTime] = useState('12:30')

  const [formData, setFormData] = useState({
    pickupName: '',
    pickupLat: '',
    pickupLon: '',
    dropoffLat: '',
    dropoffLon: '',
    tip: '',
    base: '',
    bonus: '',
    note: '',
  })

  const [errors, setErrors] = useState<Record<string, string>>({})

  function combineDateTime(date: Date, time: string): Date {
    const [hours, minutes] = time.split(':').map(Number)
    const result = new Date(date)
    result.setHours(hours, minutes, 0, 0)
    return result
  }

  function validateForm() {
    const newErrors: Record<string, string> = {}

    if (!startDate) newErrors.start = 'Start time is required'
    if (!endDate) newErrors.end = 'End time is required'
    if (!formData.pickupName) newErrors.pickupName = 'Pickup name is required'
    if (!formData.pickupLat) newErrors.pickupLat = 'Latitude is required'
    if (!formData.pickupLon) newErrors.pickupLon = 'Longitude is required'
    if (!formData.dropoffLat) newErrors.dropoffLat = 'Latitude is required'
    if (!formData.dropoffLon) newErrors.dropoffLon = 'Longitude is required'
    if (!formData.base) newErrors.base = 'Base pay is required'
    if (!formData.tip) newErrors.tip = 'Tip is required'

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()

    if (!validateForm()) return

    const start = combineDateTime(startDate, startTime)
    const end = combineDateTime(endDate, endTime)

    const request: CreateDeliveryRequest = {
      start: start.toISOString(),
      end: end.toISOString(),
      pickup: {
        name: formData.pickupName,
        lat: parseFloat(formData.pickupLat),
        lon: parseFloat(formData.pickupLon),
      },
      dropoff: {
        lat: parseFloat(formData.dropoffLat),
        lon: parseFloat(formData.dropoffLon),
      },
      earnings: {
        tip: Math.round(parseFloat(formData.tip) * 100),
        base: Math.round(parseFloat(formData.base) * 100),
        bonus: formData.bonus
          ? Math.round(parseFloat(formData.bonus) * 100)
          : undefined,
      },
      note: formData.note || undefined,
    }

    await onSubmit(request)
    setStartDate(new Date())
    setEndDate(new Date())
    setStartTime('12:00')
    setEndTime('12:30')
    setFormData({
      pickupName: '',
      pickupLat: '',
      pickupLon: '',
      dropoffLat: '',
      dropoffLon: '',
      tip: '',
      base: '',
      bonus: '',
      note: '',
    })
  }

  function updateField(field: string, value: string) {
    setFormData((prev) => ({ ...prev, [field]: value }))
    if (errors[field]) {
      setErrors((prev) => {
        const newErrors = { ...prev }
        delete newErrors[field]
        return newErrors
      })
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Create New Delivery</CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="grid grid-cols-2 gap-4">
            {/* Start Date/Time */}
            <div className="space-y-2">
              <Label>Start Time</Label>
              <Popover>
                <PopoverTrigger asChild>
                  <Button
                    variant="outline"
                    className="w-full justify-start text-left font-normal"
                    disabled={isSubmitting}
                  >
                    <CalendarIcon className="mr-2 h-4 w-4" />
                    {startDate ? (
                      `${format(startDate, 'PPP')} ${startTime}`
                    ) : (
                      <span>Pick date and time</span>
                    )}
                  </Button>
                </PopoverTrigger>
                <PopoverContent className="w-auto p-0" align="start">
                  <Calendar
                    mode="single"
                    selected={startDate}
                    onSelect={(date) => date && setStartDate(date)}
                    captionLayout="dropdown"
                  />
                  <div className="border-t p-3 space-y-2">
                    <div className="flex gap-2">
                      <Button
                        type="button"
                        variant="outline"
                        size="sm"
                        className="flex-1"
                        onClick={() => setStartDate(new Date())}
                        disabled={isSubmitting}
                      >
                        Today
                      </Button>
                      <Button
                        type="button"
                        variant="outline"
                        size="sm"
                        className="flex-1"
                        onClick={() => {
                          const yesterday = new Date()
                          yesterday.setDate(yesterday.getDate() - 1)
                          setStartDate(yesterday)
                        }}
                        disabled={isSubmitting}
                      >
                        Yesterday
                      </Button>
                    </div>
                    <Input
                      type="time"
                      value={startTime}
                      onChange={(e) => setStartTime(e.target.value)}
                      disabled={isSubmitting}
                      className="w-full"
                    />
                  </div>
                </PopoverContent>
              </Popover>
              {errors.start && (
                <p className="text-sm text-destructive">{errors.start}</p>
              )}
            </div>

            {/* End Date/Time */}
            <div className="space-y-2">
              <Label>End Time</Label>
              <Popover>
                <PopoverTrigger asChild>
                  <Button
                    variant="outline"
                    className="w-full justify-start text-left font-normal"
                    disabled={isSubmitting}
                  >
                    <CalendarIcon className="mr-2 h-4 w-4" />
                    {endDate ? (
                      `${format(endDate, 'PPP')} ${endTime}`
                    ) : (
                      <span>Pick date and time</span>
                    )}
                  </Button>
                </PopoverTrigger>
                <PopoverContent className="w-auto p-0" align="start">
                  <Calendar
                    mode="single"
                    selected={endDate}
                    onSelect={(date) => date && setEndDate(date)}
                    captionLayout="dropdown"
                  />
                  <div className="border-t p-3 space-y-2">
                    <div className="flex gap-2">
                      <Button
                        type="button"
                        variant="outline"
                        size="sm"
                        className="flex-1"
                        onClick={() => setEndDate(new Date())}
                        disabled={isSubmitting}
                      >
                        Today
                      </Button>
                      <Button
                        type="button"
                        variant="outline"
                        size="sm"
                        className="flex-1"
                        onClick={() => {
                          const yesterday = new Date()
                          yesterday.setDate(yesterday.getDate() - 1)
                          setEndDate(yesterday)
                        }}
                        disabled={isSubmitting}
                      >
                        Yesterday
                      </Button>
                    </div>
                    <Input
                      type="time"
                      value={endTime}
                      onChange={(e) => setEndTime(e.target.value)}
                      disabled={isSubmitting}
                      className="w-full"
                    />
                  </div>
                </PopoverContent>
              </Popover>
              {errors.end && (
                <p className="text-sm text-destructive">{errors.end}</p>
              )}
            </div>
          </div>

          <Separator />

          <div className="space-y-4">
            <h4 className="text-sm font-medium">Pickup Location</h4>
            <div className="space-y-2">
              <Label htmlFor="pickupName">Restaurant/Location Name</Label>
              <Input
                id="pickupName"
                placeholder="McDonald's on Main St"
                value={formData.pickupName}
                onChange={(e) => updateField('pickupName', e.target.value)}
                disabled={isSubmitting}
              />
              {errors.pickupName && (
                <p className="text-sm text-destructive">{errors.pickupName}</p>
              )}
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="pickupLat">Latitude</Label>
                <Input
                  id="pickupLat"
                  type="number"
                  step="any"
                  placeholder="40.7128"
                  value={formData.pickupLat}
                  onChange={(e) => updateField('pickupLat', e.target.value)}
                  disabled={isSubmitting}
                />
                {errors.pickupLat && (
                  <p className="text-sm text-destructive">{errors.pickupLat}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="pickupLon">Longitude</Label>
                <Input
                  id="pickupLon"
                  type="number"
                  step="any"
                  placeholder="-74.0060"
                  value={formData.pickupLon}
                  onChange={(e) => updateField('pickupLon', e.target.value)}
                  disabled={isSubmitting}
                />
                {errors.pickupLon && (
                  <p className="text-sm text-destructive">{errors.pickupLon}</p>
                )}
              </div>
            </div>
          </div>

          <Separator />

          <div className="space-y-4">
            <h4 className="text-sm font-medium">Dropoff Location</h4>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="dropoffLat">Latitude</Label>
                <Input
                  id="dropoffLat"
                  type="number"
                  step="any"
                  placeholder="40.7580"
                  value={formData.dropoffLat}
                  onChange={(e) => updateField('dropoffLat', e.target.value)}
                  disabled={isSubmitting}
                />
                {errors.dropoffLat && (
                  <p className="text-sm text-destructive">
                    {errors.dropoffLat}
                  </p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="dropoffLon">Longitude</Label>
                <Input
                  id="dropoffLon"
                  type="number"
                  step="any"
                  placeholder="-73.9855"
                  value={formData.dropoffLon}
                  onChange={(e) => updateField('dropoffLon', e.target.value)}
                  disabled={isSubmitting}
                />
                {errors.dropoffLon && (
                  <p className="text-sm text-destructive">
                    {errors.dropoffLon}
                  </p>
                )}
              </div>
            </div>
          </div>

          <Separator />

          <div className="space-y-4">
            <h4 className="text-sm font-medium">Earnings ($)</h4>
            <div className="grid grid-cols-3 gap-4">
              <div className="space-y-2">
                <Label htmlFor="base">Base Pay</Label>
                <div className="relative">
                  <span className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground text-sm">$</span>
                  <Input
                    id="base"
                    type="text"
                    inputMode="decimal"
                    placeholder="3.50"
                    value={formData.base}
                    onChange={(e) => updateField('base', e.target.value)}
                    disabled={isSubmitting}
                    className="pl-7"
                  />
                </div>
                {errors.base && (
                  <p className="text-sm text-destructive">{errors.base}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="tip">Tip</Label>
                <div className="relative">
                  <span className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground text-sm">$</span>
                  <Input
                    id="tip"
                    type="text"
                    inputMode="decimal"
                    placeholder="2.00"
                    value={formData.tip}
                    onChange={(e) => updateField('tip', e.target.value)}
                    disabled={isSubmitting}
                    className="pl-7"
                  />
                </div>
                {errors.tip && (
                  <p className="text-sm text-destructive">{errors.tip}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="bonus">Bonus (Optional)</Label>
                <div className="relative">
                  <span className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground text-sm">$</span>
                  <Input
                    id="bonus"
                    type="text"
                    inputMode="decimal"
                    placeholder="1.50"
                    value={formData.bonus}
                    onChange={(e) => updateField('bonus', e.target.value)}
                    disabled={isSubmitting}
                    className="pl-7"
                  />
                </div>
              </div>
            </div>
          </div>

          <Separator />

          <div className="space-y-2">
            <Label htmlFor="note">Note (Optional)</Label>
            <Textarea
              id="note"
              placeholder="Add any notes about this delivery..."
              value={formData.note}
              onChange={(e) => updateField('note', e.target.value)}
              disabled={isSubmitting}
            />
          </div>

          <Button type="submit" className="w-full" disabled={isSubmitting}>
            {isSubmitting ? (
              <>
                <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                Creating...
              </>
            ) : (
              'Create Delivery'
            )}
          </Button>
        </form>
      </CardContent>
    </Card>
  )
}
