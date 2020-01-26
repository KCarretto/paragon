import * as React from "react";
import { FunctionComponent, useState } from "react";
import { Popup } from "semantic-ui-react";

type SidebarProps = {
  value: string;
};

const XClipboard: FunctionComponent<SidebarProps> = props => {
  const [copySuccess, setCopySuccess] = useState<string>("Copy to clipboard!");
  let copy = (_: React.MouseEvent<HTMLSpanElement, MouseEvent>) => {
    navigator.clipboard.writeText(props.value);
    setCopySuccess("Successfully copied '" + props.value + "' to clipboard!");
    setTimeout(function() {
      setCopySuccess("Copy to clipboard!");
    }, 1500);
  };
  return (
    <Popup
      hoverable
      wide
      closeOnTriggerClick={false}
      popperDependencies={[copySuccess]}
      size="mini"
      content={copySuccess}
      trigger={
        <span style={{ cursor: "pointer" }} onClick={copy}>
          {props.children}
        </span>
      }
    />
  );
};
export default XClipboard;
