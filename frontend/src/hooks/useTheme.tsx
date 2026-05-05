import { createContext, useContext, useEffect, useState } from "react";

// All available themes
export type Theme =
  // Catppuccin flavors
  | "ctp-latte"
  | "ctp-frappe"
  | "ctp-macchiato"
  | "ctp-mocha"
  // TokyoNight variants
  | "tokyonight-day"
  | "tokyonight-moon"
  | "tokyonight-night"
  | "tokyonight-storm";

interface ThemeContextType {
  currentTheme: Theme;
  setTheme: (theme: Theme) => void;
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined);

// Theme metadata for UI display
export const themeMetadata: Record<
  Theme,
  { label: string; family: string; isDark: boolean }
> = {
  // Catppuccin
  "ctp-latte": {
    label: "Catppuccin Latte",
    family: "catppuccin",
    isDark: false,
  },
  "ctp-frappe": {
    label: "Catppuccin Frappé",
    family: "catppuccin",
    isDark: true,
  },
  "ctp-macchiato": {
    label: "Catppuccin Macchiato",
    family: "catppuccin",
    isDark: true,
  },
  "ctp-mocha": {
    label: "Catppuccin Mocha",
    family: "catppuccin",
    isDark: true,
  },
  // TokyoNight
  "tokyonight-day": {
    label: "TokyoNight Day",
    family: "tokyonight",
    isDark: false,
  },
  "tokyonight-moon": {
    label: "TokyoNight Moon",
    family: "tokyonight",
    isDark: true,
  },
  "tokyonight-night": {
    label: "TokyoNight Night",
    family: "tokyonight",
    isDark: true,
  },
  "tokyonight-storm": {
    label: "TokyoNight Storm",
    family: "tokyonight",
    isDark: true,
  },
};

// All available themes as an array
export const availableThemes: Theme[] = Object.keys(themeMetadata) as Theme[];

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const [currentTheme, setCurrentTheme] = useState<Theme>(() => {
    const savedTheme = localStorage.getItem("theme") as Theme | null;
    return savedTheme && savedTheme in themeMetadata
      ? savedTheme
      : "tokyonight-night";
  });

  useEffect(() => {
    // Remove all theme classes first
    const allThemes = availableThemes;
    document.documentElement.classList.remove(...allThemes);

    // Remove latte class (will be re-added if needed)
    document.documentElement.classList.remove("latte");

    // Add current theme class
    document.documentElement.classList.add(currentTheme);

    // Add latte class for Catppuccin Latte theme to enable light mode variant
    if (currentTheme === "ctp-latte") {
      document.documentElement.classList.add("latte");
    }

    // Save to localStorage
    localStorage.setItem("theme", currentTheme);
  }, [currentTheme]);

  return (
    <ThemeContext.Provider value={{ currentTheme, setTheme: setCurrentTheme }}>
      {children}
    </ThemeContext.Provider>
  );
}

export function useTheme() {
  const context = useContext(ThemeContext);
  if (!context) {
    throw new Error("useTheme must be used within ThemeProvider");
  }
  return context;
}
