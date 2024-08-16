'use client';

import React, { useRef, useState } from 'react';
import { Editor } from '@monaco-editor/react';

const Test = () => {
  const editorRef = useRef();
  const [value,setValue] = useState('// write your code here');

  const onMount = (editor) => {
    editorRef.current = editor;
    editor.focus();
  }
  return (
    <main className='min-h-screen bg-white dark:bg-black text-black dark:text-white flex flex-col content-center'>
      <Editor 
        height="90vh" 
        width="50vw"
        defaultLanguage="python" 
        value={value}
        onChange={(value) => setValue(value)}
        onMount={onMount}
        theme='vs-dark'
        saveViewState={true}
        className='border border-black dark:border-white mt-20'
      />
      <div className='w-1/2 flex justify-end p-2 h-max bg-white dark:bg-black text-black dark:text-white'>
        <button type='button' className='rounded-sm border-white border'>RUN</button>
      </div>
    </main>
  )
}

export default Test