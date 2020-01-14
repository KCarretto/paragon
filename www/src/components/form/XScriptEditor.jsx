import { ControlledEditor } from "@monaco-editor/react";
import React from 'react';
import { XLoadingMessage } from '../messages';

const XScriptEditor = ({ onChange, content }) => {
    return (
        <ControlledEditor
            loading={<XLoadingMessage
                title='Loading Script Editor'
                msg='Initializing scripting engine...'
            />}
            height='50vh'
            language="python"
            theme='dark'
            value={content}
            onChange={(e, value) => onChange(e, { value: value })}
            options={{
                scrollbar: {
                    verticalScrollbarSize: '7px',
                },
                minimap: { enabled: false },
                cursorStyle: 'line-thin',
            }}
            editorDidMount={(fn, mco) => {
                let element = document.getElementsByTagName('textarea')[0];
                element.classList.remove("inputarea");
            }}
        />
    );
}
export default XScriptEditor;