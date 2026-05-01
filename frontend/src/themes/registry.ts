import { flavors } from '@catppuccin/palette'

const catppuccinMocha = flavors.mocha

export interface ThemeDefinition {
  name: string
  label: string
  colors: Record<string, string>
}

/**
 * Centralized theme registry.
 * Imports from official packages where available (@catppuccin/palette).
 * Falls back to official source definitions for Tokyo Night and Gruvbox.
 */
export const themeRegistry: Record<string, ThemeDefinition> = {
  catppuccin: {
    name: 'catppuccin',
    label: 'Catppuccin Mocha',
    colors: {
      '--color-bg-primary': catppuccinMocha.colors.base.hex,
      '--color-bg-secondary': catppuccinMocha.colors.mantle.hex,
      '--color-bg-tertiary': catppuccinMocha.colors.surface0.hex,
      '--color-text-primary': catppuccinMocha.colors.text.hex,
      '--color-text-secondary': catppuccinMocha.colors.subtext0.hex,
      '--color-accent': catppuccinMocha.colors.blue.hex,
      '--color-accent-hover': catppuccinMocha.colors.sky.hex,
      '--color-border': catppuccinMocha.colors.overlay0.hex,
      '--color-success': catppuccinMocha.colors.green.hex,
      '--color-error': catppuccinMocha.colors.red.hex,
    },
  },
  'tokyo-night': {
    name: 'tokyo-night',
    label: 'Tokyo Night',
    colors: {
      // Official palette: https://github.com/tokyo-night/tokyo-night-vscode-theme
      '--color-bg-primary': '#1a1b26',
      '--color-bg-secondary': '#16161e',
      '--color-bg-tertiary': '#24283b',
      '--color-text-primary': '#c0caf5',
      '--color-text-secondary': '#565f89',
      '--color-accent': '#7aa2f7',
      '--color-accent-hover': '#bb9af7',
      '--color-border': '#3b4261',
      '--color-success': '#9ece6a',
      '--color-error': '#f7768e',
    },
  },
  gruvbox: {
    name: 'gruvbox',
    label: 'Gruvbox Dark (Medium)',
    colors: {
      // Official palette: https://github.com/morhetz/gruvbox
      '--color-bg-primary': '#282828',
      '--color-bg-secondary': '#1d2021',
      '--color-bg-tertiary': '#3c3836',
      '--color-text-primary': '#ebdbb2',
      '--color-text-secondary': '#a89984',
      '--color-accent': '#83a598',
      '--color-accent-hover': '#b8bb26',
      '--color-border': '#504945',
      '--color-success': '#b8bb26',
      '--color-error': '#fb4934',
    },
  },
}
