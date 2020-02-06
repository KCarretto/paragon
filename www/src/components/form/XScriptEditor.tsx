import { ControlledEditor } from "@monaco-editor/react";
import React from "react";
import { XLoadingMessage } from "../messages";

// monaco.init().then(monaco => {
//   Register(monaco);
// });

// monaco.init().then(monaco => {
//   monaco.languages.register({ id: "renegade" });
//   monaco.languages.setMonarchTokensProvider("renegade", renegade_language);
//   monaco.languages.setLanguageConfiguration("renegade", renegade_conf(monaco));
//   monaco.languages.registerCompletionItemProvider(
//     "renegade",
//     renegade_autocomplete(monaco)
//   );
// });

const XScriptEditor = ({ onChange, content }) => {
  return (
    <ControlledEditor
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
  );
};
export default XScriptEditor;
