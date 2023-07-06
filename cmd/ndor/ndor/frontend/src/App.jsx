import { useState } from 'react';
import MonacoEditor from 'react-monaco-editor';

import './App.css';

import { GenImage } from "../wailsjs/go/main/APP"
import { AlertError } from '../wailsjs/go/main/APP';

function App() {
  const [genedImgSrc, setGenedImgSrc] = useState('')
  const [isEditorReady, setIsEditorReady] = useState(false);
  const [showExamples, setShowExamples] = useState(false)
  const [code, setCode] = useState('')
  const [editor, setEditor] = useState()

  function handleEditorDidMount(editor, monaco) {
    setIsEditorReady(true);
    setEditor(editor)
    editor.focus()
  }
  function onEditorChange(newValue, e) {
    setCode(newValue)
  }

  function help() {

  }
  function go() {
    console.log('code to run:', code)
    GenImage(code).then((res) => {
      const src = res['src']
      const err = res['err']
      if (src != '') {
        setGenedImgSrc(src)
      } else {
        AlertError(err)
      }
    })
  }

  return (
    <div style={{
      'overflow': "hidden", "--wails-draggable": "drag"
    }} >
      < div className='title-bar' >
        <p>Ndor 牛豆儿画图</p>
      </div>
      <div id="pictureBox" className='left-box'>
        <img id='pictureArea' src={genedImgSrc} style={{ 'max- width': "100%", 'max-height': "100%" }} />
      </ div >
      <div className='right-box'>
        <MonacoEditor
          language="javascript"
          theme="vs-dark"
          options={{
            fontSize: 16,
            wordWrap: true,
            minimap: {
              enabled: false
            },
            scrollBar: {
              vertical: 'hidden',
              horizontal: 'hidden'
            },
            automaticLayout: true,
            overviewrulerLanes: 0,
            hideCursorInOverviewRuler: true,
          }}
          editorDidMount={handleEditorDidMount}
          onChange={onEditorChange}
        />
      </div>
      <button className='teacher-button' onClick={help()}>HELP</button>
      <button id="run-button" className='run-button' disabled={!isEditorReady} onClick={go}>GO</button>
    </ div >
  )

}
export default App
