import { ControlledEditor } from "@monaco-editor/react";
import React, { useState } from 'react';

const XScriptEditor = ({ handleChange, defaultContent }) => {
    if (!defaultContent) {
        defaultContent = `\n# Enter your script here!\n\ndef main():\n\tprint("Hello World!")`;
    }
    const [content, setContent] = useState(defaultContent);

    return (
        <ControlledEditor
            options={{
                scrollbar: {
                    verticalScrollbarSize: '7px',
                },
                minimap: { enabled: false },
                cursorStyle: 'line-thin',
            }}
            theme='dark'
            height='550px'
            value={content}
            editorDidMount={(fn, mco) => {
                let element = document.getElementsByTagName('textarea')[0];
                element.classList.remove("inputarea");
            }}
            onChange={(e, value) => {
                setContent(value);
                handleChange(e, { name: 'content', value: value });
            }}
            language="python"
        />
    );
}
export default XScriptEditor;