import { writeFileSync } from 'fs'
import { themeRegistry } from '../src/themes/registry'

let cssOutput = '/* Auto-generated from src/themes/registry.ts — do not edit manually */\n\n'

for (const [themeKey, themeDef] of Object.entries(themeRegistry)) {
  const selector = themeKey === 'catppuccin' ? ':root' : `[data-theme="${themeKey}"]`
  
  cssOutput += `${selector} {\n`
  for (const [cssVariable, colorValue] of Object.entries(themeDef.colors)) {
    cssOutput += `  ${cssVariable}: ${colorValue};\n`
  }
  cssOutput += '}\n\n'
}

writeFileSync('src/styles/themes.css', cssOutput)
console.log('✅ themes.css generated successfully')
