import { ControlledEditor } from "@monaco-editor/react";
import React from 'react';

const XScriptEditor = ({ onChange, content }) => {
    console.log("SCRIPT EDITOR RENDERED")
    // const name = 'content';
    // const [state, setState] = useState({ content: `\n# Enter your script here!\n\ndef main():\n\tprint("Hello World!")` });

    // const handleChange = (e, content) => {
    //     console.log("SCRIPT EDITOR EVENT: ", e)
    //     onChange(e, { name: 'content', value: content });
    //     return content;
    // }

    return (
        // <Form.Field style={{ 'margin-top': '25px' }}>
        //     <label>Script</label>
        //     <Form.TextArea
        //         label={{ content: 'Enter script' }}
        //         placeholder='Enter script content'
        //         name='content'
        //         rows={15}
        //         value={state.content}
        //         onChange={(e, { value }) => handleChange(e, value)}
        //     />
        // </Form.Field>
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
            onChange={(e, value) => onChange(e, { value: value })}
            language="python"
        />
    );
}
export default XScriptEditor;