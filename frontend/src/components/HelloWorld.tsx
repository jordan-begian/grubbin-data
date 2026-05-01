import { useHelloWorldGreeting } from '../hooks/useHelloWorld'
import styles from './HelloWorld.module.css'

export default function HelloWorld() {
  const { data: greetingData, isLoading, isError, error } = useHelloWorldGreeting()

  if (isLoading) {
    return <div className={styles.loadingContainer}>Loading greeting...</div>
  }

  if (isError) {
    return <div className={styles.errorContainer}>Error: {error.message}</div>
  }

  return (
    <article className={styles.greetingContainer}>
      <h1 className={styles.greetingMessage}>{greetingData?.message}</h1>
      <time className={styles.generatedAtTimestamp} dateTime={greetingData?.timestamp}>
        {greetingData?.timestamp}
      </time>
    </article>
  )
}
