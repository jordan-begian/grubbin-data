import { useTheme, themeMetadata, availableThemes } from '../hooks/useTheme'

export function ThemeSwitcher() {
  const { currentTheme, setTheme } = useTheme()

  return (
    <div 
      className="flex items-center gap-2 px-8 py-4 max-w-3xl mx-auto"
      style={{ backgroundColor: 'var(--color-bg-primary)' }}
    >
      <label 
        htmlFor="theme-select" 
        className="text-sm font-medium"
        style={{ color: 'var(--color-text-secondary)' }}
      >
        Theme:
      </label>
      <select
        id="theme-select"
        value={currentTheme}
        onChange={(event) => setTheme(event.target.value as typeof currentTheme)}
        className="rounded-md px-3 py-1.5 text-sm cursor-pointer transition-colors focus:outline-none focus:ring-2"
        style={{
          backgroundColor: 'var(--color-bg-tertiary)',
          color: 'var(--color-text-primary)',
          border: '1px solid var(--color-border)',
        }}
      >
        <optgroup label="Catppuccin">
          {availableThemes
            .filter((theme) => themeMetadata[theme].family === 'catppuccin')
            .map((theme) => (
              <option key={theme} value={theme}>
                {themeMetadata[theme].label}
              </option>
            ))}
        </optgroup>
        <optgroup label="TokyoNight">
          {availableThemes
            .filter((theme) => themeMetadata[theme].family === 'tokyonight')
            .map((theme) => (
              <option key={theme} value={theme}>
                {themeMetadata[theme].label}
              </option>
            ))}
        </optgroup>
      </select>
    </div>
  )
}
