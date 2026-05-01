import { useTheme } from '../hooks/useTheme'
import { themeRegistry } from '../themes/registry'
import styles from './ThemeSwitcher.module.css'

export function ThemeSwitcher() {
  const { currentTheme, setTheme } = useTheme()

  return (
    <div className={styles.themeSwitcherContainer}>
      <label htmlFor="theme-select" className={styles.themeLabel}>
        Theme:
      </label>
      <select
        id="theme-select"
        className={styles.themeSelect}
        value={currentTheme}
        onChange={(event) => setTheme(event.target.value as typeof currentTheme)}
      >
        {Object.values(themeRegistry).map((themeOption) => (
          <option key={themeOption.name} value={themeOption.name}>
            {themeOption.label}
          </option>
        ))}
      </select>
    </div>
  )
}
