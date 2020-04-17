import * as monaco from "monaco-editor";
import * as React from "react";
import { useEffect } from "react";
import { Icon, Menu } from "semantic-ui-react";
import * as Renegade from "../../config/renegade";

Renegade.Register();

export interface XEditorProps {
  value: string;
  onChange: (e: monaco.editor.IModelContentChangedEvent, value: string) => void;
}

const XEditor: React.FunctionComponent<XEditorProps> = ({
  value,
  onChange
}: XEditorProps) => {
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
      language: Renegade.LanguageID,
      value: value,
      scrollbar: {
        verticalScrollbarSize: 7
      },
      fixedOverflowWidgets: true,
      minimap: {
        enabled: false
      },
      cursorStyle: "line-thin"
    });

    element.getElementsByTagName("textarea")[0].classList.remove("inputarea");
    editor.onDidChangeModelContent(e => {
      onChange(e, editor.getValue());
      let action = editor.getAction("editor.action.triggerParameterHints");
      if (action !== null) {
        action.run().then(v => console.log("RAN ACTION", v));
      }
    });
  }, []);

  return (
    <React.Fragment>
      {/* <div className="XEditorFileMenu"> */}
      <Menu
        compact
        inverted
        attached="top"
        style={{
          marginTop: "0px",
          borderRadius: "0px",
          backgroundColor: "rgba(63, 63, 63, 1)"
        }}
      >
        <Menu.Item
          active
          fitted="vertically"
          onClick={e => {
            e.preventDefault();
            e.stopPropagation();
            alert("OPEN FILE");
          }}
        >
          <span className="XEditorFileMenuItemLabel">main.rg</span>
          <Icon
            name="delete"
            onClick={e => {
              e.preventDefault();
              e.stopPropagation();
              alert("CLOSE FILE");
            }}
          />
        </Menu.Item>
        <Menu.Item
          fitted="vertically"
          onClick={e => {
            e.preventDefault();
            e.stopPropagation();
            alert("OPEN FILE");
          }}
        >
          <span className="XEditorFileMenuItemLabel">deploy_nomnom.rg</span>
          <Icon
            name="delete"
            onClick={e => {
              e.preventDefault();
              e.stopPropagation();
              alert("CLOSE FILE");
            }}
          />
        </Menu.Item>
        <Menu.Item
          fitted="vertically"
          onClick={e => {
            e.preventDefault();
            e.stopPropagation();
            alert("OPEN FILE");
          }}
        >
          <span className="XEditorFileMenuItemLabel">deploy_paragon.rg</span>
          <Icon
            name="delete"
            onClick={e => {
              e.preventDefault();
              e.stopPropagation();
              alert("CLOSE FILE");
            }}
          />
        </Menu.Item>
      </Menu>
      {/* </div> */}
      {container}
    </React.Fragment>
  );
};

export default XEditor;
