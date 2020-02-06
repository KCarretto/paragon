import * as monaco from "monaco-editor";
import * as React from "react";
import { useEffect } from "react";

const XEditor = () => {
  const container = <div id="editor" className="XEditor" />;

  /*
loading={
        <XLoadingMessage
          title="Loading Script Editor"
          msg="Initializing scripting engine..."
        />
      }
      height="50vh"
      language="renegade"
      theme="renegade"
      value={content}
      onChange={(e, value) => onChange(e, { value: value })}
      options={{
        scrollbar: {
          verticalScrollbarSize: 7
        },
        minimap: { enabled: false },
        cursorStyle: "line-thin"
      }}
      editorDidMount={(fn, mco) => {
        let element = document.getElementsByTagName("textarea")[0];
        element.classList.remove("inputarea");
      }}
    />
*/
  useEffect(() => {
    let element = document.getElementById("editor");
    let editor = monaco.editor.create(element, {
      theme: "renegade",
      language: "renegade",
      scrollbar: {
        verticalScrollbarSize: 7
      },
      minimap: {
        enabled: false
      },
      cursorStyle: "line-thin"
    });
  }, []);

  return container;
};

export default XEditor;
