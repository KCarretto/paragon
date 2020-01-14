import { ControlledEditor } from "@monaco-editor/react";
import React from 'react';
import { Icon, Message } from 'semantic-ui-react';

const XScriptEditor = ({ onChange, content }) => {
    return (
        <ControlledEditor
            loading={<Message icon size='massive'>
                <Icon name='circle notched' loading />
                <Message.Content>
                    <Message.Header>Loading Script Editor</Message.Header>
                    Initializing scripting engine...
                </Message.Content>
            </Message>}
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