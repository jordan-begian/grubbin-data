import { useHelloWorldGreeting } from '../hooks/useHelloWorld'

export default function HelloWorld() {
  const { data: greetingData, isLoading, isError, error } = useHelloWorldGreeting()

  if (isLoading) {
    return (
      <div 
        className="text-center mt-16"
        style={{ color: 'var(--color-text-secondary)' }}
      >
        Loading greeting...
      </div>
    )
  }

  if (isError) {
    return (
      <div 
        className="text-center mt-16"
        style={{ color: '#f38ba8' }}  /* Catppuccin red as fallback */
      >
        Error: {error.message}
      </div>
    )
  }

  return (
    <article 
      className="max-w-xl mx-auto my-16 p-8 text-center border rounded-lg"
      style={{
        backgroundColor: 'var(--color-bg-secondary)',
        borderColor: 'var(--color-border)',
      }}
    >
      <h1 
        className="text-3xl mb-4"
        style={{ color: 'var(--color-text-primary)' }}
      >
        {greetingData?.message}
      </h1>
      <time 
        className="text-sm"
        style={{ color: 'var(--color-text-secondary)' }}
        dateTime={greetingData?.timestamp}
      >
        {greetingData?.timestamp}
      </time>
    </article>
  )
}
