import { BrowserRouter } from 'react-router-dom'
import HelloWorld from './components/HelloWorld'
import { ThemeSwitcher } from './components/ThemeSwitcher'

export default function App() {
  return (
    <BrowserRouter>
      <main>
        <ThemeSwitcher />
        <HelloWorld />
      </main>
    </BrowserRouter>
  )
}
