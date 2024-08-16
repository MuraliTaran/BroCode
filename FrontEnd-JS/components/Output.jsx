'use client';

import React, { useState } from 'react';
import { AXIOS, VERSIONS, BASE_URL } from '@/utils/constants';

const Output = ({value, lang}) => {
  const [showSubs, setShowSubs] = useState(false);
  const [logs, setLogs] = useState([]);
  const [submissions, setSubmissions] = useState([]);
  const Run_Code = async (lang, code) => {
    try{
      const response = await AXIOS.post(BASE_URL + '/run', 
      {
        "language" : lang,
        "version" : VERSIONS[lang],
        "code" : code,
        "stdin" : "3\n1\n2\n3"
      })
      if (response.data.status) {
        setShowSubs(false)
        setLogs([...logs, response.data.output])
      }
    }
    catch(error) {
      console.log("error", error)
    }
  }
  const Submit_Code = async (lang, code) => {
    try{
      const response = await AXIOS.post(BASE_URL + '/submit', 
      {
        "user_id" : "65df3889ed28f385c98cb76c",
        "problem_id" : "6601dd07194df4b680303eb8",
        "language" : lang,
        "version" : VERSIONS[lang],
        "code" : code,
      })
      if (response.data.status) {
        setShowSubs(true)
        setSubmissions([...submissions, response.data])
      }
    }
    catch(error) {
      console.log("error", error)
    }
  }
  return (
    <section className='w-1/2 flex flex-col px-1 mb-2'>
      <div className='flex flex-row gap-2'>
        <button className={`w-32 p-2 rounded-sm my-2 border border-white bg-black`} onClick={() => Run_Code(lang, value)}>RUN</button>
        <button className={`w-32 p-2 rounded-sm my-2 border border-white bg-black`} onClick={() => Submit_Code(lang, value)}>SUBMIT</button>
      </div>
      <div className='flex-1 bg-neutral-900 border-white border rounded-sm overflow-y-scroll '>
        {showSubs ? 
        <>
        {submissions.map((sub) => (
          <div className='rounded-sm px-2 flex flex-col border-b-[0.5px] border-gray-500 last:bg-slate-800'>
            <p>{`result : ${sub.code}`}</p>
            <p>{`passed : ${sub.passed_cases}`}</p>
            <p>{`total  : ${sub.total_cases}`}</p>
            {sub.one_of_failed &&
            <>
              <p>{`one of failed case : ${sub.one_of_failed.test_case}`}</p>
              <p>{`expected          : ${sub.one_of_failed.expected}`}</p>
              <p>{`output            : ${sub.one_of_failed.output}`}</p>
            </>
            }
          </div>
        ))}
        </> :
        <>
        {logs.map((log) => (
          <div className='rounded-sm px-2 flex flex-row border-b-[0.5px] border-gray-500 last:bg-slate-800'>
            <div className='mr-1'>{`>>`}</div>
            <div>
            {log.replace(/\n+$/, '').split('\n').map((text, index) => (
            <React.Fragment key={index}>
              {text}
              <br />
            </React.Fragment>
            ))}
            </div>
          </div>
        ))}
        </>
        }
      </div>
    </section>
  )
}

export default Output