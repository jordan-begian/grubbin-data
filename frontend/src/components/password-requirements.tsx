import { cn } from "@/lib/utils"
import { validatePassword, type PasswordValidationResult } from "@/lib/password-validation"

interface PasswordRequirementsProps {
  password: string
  className?: string
}

function RequirementItem({ label, met }: { label: string; met: boolean }) {
  return (
    <li
      className={cn(
        "flex items-center gap-2 text-sm transition-colors duration-200",
        met ? "text-green-600 dark:text-green-400" : "text-muted-foreground"
      )}
    >
      <span
        className={cn(
          "flex h-4 w-4 items-center justify-center rounded-full text-xs font-bold",
          met
            ? "bg-green-100 text-green-600 dark:bg-green-900/30 dark:text-green-400"
            : "bg-muted text-muted-foreground"
        )}
        aria-hidden="true"
      >
        {met ? "\u2713" : "\u2717"}
      </span>
      {label}
    </li>
  )
}

export function PasswordRequirements({ password, className }: PasswordRequirementsProps) {
  const result: PasswordValidationResult = validatePassword(password)
  const allMet = result.isValid

  if (password.length === 0) {
    return null
  }

  return (
    <div className={cn("space-y-2", className)}>
      <p
        className={cn(
          "text-sm font-medium transition-colors duration-200",
          allMet ? "text-green-600 dark:text-green-400" : "text-muted-foreground"
        )}
      >
        {allMet ? "All requirements met" : "Password requirements:"}
      </p>
      <ul className="space-y-1">
        {result.requirements.map((req) => (
          <RequirementItem key={req.label} label={req.label} met={req.met} />
        ))}
      </ul>
    </div>
  )
}
