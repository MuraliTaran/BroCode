'use client';

import React, { useRef, useState } from 'react';
import { Editor } from '@monaco-editor/react';
import Output from './Output';
import { SNIPPETS } from '@/utils/constants';

const CodePlace = () => {
  const editorRef = useRef();
  const [lang, setLang] = useState('go');
  const [value, setValue] = useState(SNIPPETS[lang]);
  const OnMount = (editor) => {
    editorRef.current = editor
    editor.focus()
  }

  return (
    <section className='flex flex-row w-full gap-1'>
      <div className='flex-1 flex flex-col '>
        <div className='flex flex-row gap-2 my-2'>
          <button className={`w-32 p-2 rounded-sm border border-white ${lang === "go" ? `bg-stone-800` :`bg-black`}`} onClick={() => {setLang('go'); setValue(SNIPPETS['go'])}}>GO</button>
          <button className={`w-32 p-2 rounded-sm border border-white ${lang === "python" ? `bg-stone-800` :`bg-black`}`} onClick={() => {setLang('python'); setValue(SNIPPETS['python'])}}>PYTHON</button>
        </div>
        <Editor 
          height='90vh'
          theme='vs-dark'
          language={lang}
          defaultValue={SNIPPETS[lang]}
          value={value}
          onChange={(value) => setValue(value)}
          onMount={OnMount}
          options={{minimap: {enabled: false}, 
                    contextmenu: false
                  }}
          className='border-white border rounded-sm'
        />
      </div>
      <Output value={value} lang={lang} />
    </section>
  )
}

export default CodePlace