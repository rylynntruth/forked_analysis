import type { GetServerSideProps } from 'next'
import Page from '../components/page'
import { initializeStore } from '../lib/store'

export default function Home() {
  return <Page />
}

// The date returned here will be different for every request that hits the page,
// that is because the page becomes a serverless function instead of being statically
// exported when you use `getServerSideProps` or `getInitialProps`
export const getServerSideProps: GetServerSideProps = async () => {
  const zustandStore = initializeStore()

  return {
    props: {
      // the "stringify and then parse again" piece is required as next.js
      // isn't able to serialize it to JSON properly
      initialZustandState: JSON.parse(JSON.stringify(zustandStore.getState())),
    },
  }
}
