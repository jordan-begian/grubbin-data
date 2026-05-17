
import { useState } from 'react'
import { format } from 'date-fns'
import { Calendar as CalendarIcon } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Separator } from '@/components/ui/separator'
import { Calendar } from '@/components/ui/calendar'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'
import { Loader2 } from 'lucide-react'
import type { DeliveryResponse, UpdateDeliveryRequest } from '@/types/delivery'
import { useForm } from 'react-hook-form'

interface DeliveryEditFormProps {
  delivery: DeliveryResponse
  onSubmit: (data: UpdateDeliveryRequest) => Promise<void>
  isSubmitting: boolean
}

interface EditFormData {
  pickupName: string
  base: string
  tip: string
  bonus: string
  note: string
}

function parseDateTime(isoString: string): { date: Date; time: string } {
  const date = new Date(isoString)
  const hours = date.getHours().toString().padStart(2, '0')
  const minutes = date.getMinutes().toString().padStart(2, '0')
  return { date, time: `${hours}:${minutes}` }
}

function combineDateTime(date: Date, time: string): string {
  const [hours, minutes] = time.split(':').map(Number)
  const result = new Date(date)
  result.setHours(hours, minutes, 0, 0)
  return result.toISOString()
}

export function DeliveryEditForm({
  delivery,
  onSubmit,
  isSubmitting,
}: DeliveryEditFormProps) {
  const startParsed = parseDateTime(delivery.start)
  const endParsed = parseDateTime(delivery.end)

  const [startDate, setStartDate] = useState<Date>(startParsed.date)
  const [startTime, setStartTime] = useState(startParsed.time)
  const [endDate, setEndDate] = useState<Date>(endParsed.date)
  const [endTime, setEndTime] = useState(endParsed.time)

  const form = useForm<EditFormData>({
    defaultValues: {
      pickupName: delivery.pickup.name,
      base: (delivery.earnings.base / 100).toString(),
      tip: (delivery.earnings.tip / 100).toString(),
      bonus: delivery.earnings.bonus
        ? (delivery.earnings.bonus / 100).toString()
        : '',
      note: delivery.note || '',
    },
  })

  async function handleSubmit(data: EditFormData) {
    const update: UpdateDeliveryRequest = {
      id: delivery.id,
    }

    const newStart = combineDateTime(startDate, startTime)
    const newEnd = combineDateTime(endDate, endTime)

    if (newStart !== delivery.start) {
      update.start = newStart
    }
    if (newEnd !== delivery.end) {
      update.end = newEnd
    }
    if (data.pickupName !== delivery.pickup.name) {
      update.pickup = {
        name: data.pickupName,
        lat: delivery.pickup.lat,
        lon: delivery.pickup.lon,
      }
    }

    const baseCents = Math.round(parseFloat(data.base) * 100)
    const tipCents = Math.round(parseFloat(data.tip) * 100)
    const bonusCents = data.bonus
      ? Math.round(parseFloat(data.bonus) * 100)
      : undefined

    if (
      baseCents !== delivery.earnings.base ||
      tipCents !== delivery.earnings.tip ||
      bonusCents !== delivery.earnings.bonus
    ) {
      update.earnings = {
        tip: tipCents,
        base: baseCents,
        bonus: bonusCents,
      }
    }

    if (data.note !== (delivery.note || '')) {
      update.note = data.note || null
    }

    await onSubmit(update)
  }

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(handleSubmit)}
        className="space-y-6"
      >
        <div className="grid grid-cols-2 gap-4">
          {/* Start Date/Time */}
          <FormItem>
            <FormLabel>Start Time</FormLabel>
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
          </FormItem>

          {/* End Date/Time */}
          <FormItem>
            <FormLabel>End Time</FormLabel>
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
          </FormItem>
        </div>

        <Separator />

        <div className="space-y-4">
          <h4 className="text-sm font-medium">Pickup Location</h4>
          <FormField
            control={form.control}
            name="pickupName"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Name</FormLabel>
                <FormControl>
                  <Input {...field} disabled={isSubmitting} />
                </FormControl>
              </FormItem>
            )}
          />
        </div>

        <Separator />

        <div className="space-y-4">
          <h4 className="text-sm font-medium">Earnings ($)</h4>
          <div className="grid grid-cols-3 gap-4">
            <FormField
              control={form.control}
              name="base"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Base Pay</FormLabel>
                  <FormControl>
                    <Input
                      type="number"
                      step="0.01"
                      {...field}
                      disabled={isSubmitting}
                    />
                  </FormControl>
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="tip"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Tip</FormLabel>
                  <FormControl>
                    <Input
                      type="number"
                      step="0.01"
                      {...field}
                      disabled={isSubmitting}
                    />
                  </FormControl>
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="bonus"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Bonus</FormLabel>
                  <FormControl>
                    <Input
                      type="number"
                      step="0.01"
                      {...field}
                      value={field.value || ''}
                      disabled={isSubmitting}
                    />
                  </FormControl>
                </FormItem>
              )}
            />
          </div>
        </div>

        <Separator />

        <FormField
          control={form.control}
          name="note"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Note</FormLabel>
              <FormControl>
                <Textarea {...field} disabled={isSubmitting} />
              </FormControl>
            </FormItem>
          )}
        />

        <Button type="submit" className="w-full" disabled={isSubmitting}>
          {isSubmitting ? (
            <>
              <Loader2 className="h-4 w-4 mr-2 animate-spin" />
              Saving...
            </>
          ) : (
            'Save Changes'
          )}
        </Button>
      </form>
    </Form>
  )
}
