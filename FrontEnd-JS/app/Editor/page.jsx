import React from 'react';
import CodePlace from '@/components/codeplace';
import Output from '@/components/Output';

const page = () => {
  return (
    <main className='h-screen flex flex-row bg-white dark:bg-black dark:text-white text-black'>
      <CodePlace className='flex-1'/>
    </main>
  )
}

export default page