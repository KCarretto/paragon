import * as React from "react";
import { Button, Icon, Label } from "semantic-ui-react";

type XFileInputProps = {
  id: string;
  setFile: (file: File | null) => void;
  file: File | null;
};

const XFileInput: React.FunctionComponent<XFileInputProps> = ({
  id,
  setFile,
  file,
}) => {
  let inputRef: HTMLInputElement | null = null;

  return (
    <React.Fragment>
      {file === null ? (
        <Button
          icon="upload"
          htmlFor={id}
          label={
            <Label
              as="label"
              style={{ cursor: "pointer" }}
              basic
              children="Select file"
            />
          }
          onClick={() => inputRef!.click()}
          labelPosition="right"
        />
      ) : (
          <Label>
            <Icon name="file alternate outline" />
            {file.name} ({file.size} bytes)
            <Icon
              name="delete"
              onClick={() => {
                inputRef.value = null;
                setFile(null);
              }}
            />
          </Label>
        )}
      <input
        hidden
        ref={el => {
          inputRef = el!;
        }}
        type="file"
        id={id}
        onChange={e => {
          e.preventDefault();
          console.log("FILE UPLOAD EVENT", e);
          console.log("FILE UPLOAD EVENT TARGET", e.target);
          console.log("FILE UPLOAD EVENT TARGET VALUE", e.target.value);
          console.log("FILE UPLOAD EVENT TARGET FILES", e.target.files);
          setFile(e.target.files[0]);
        }}
      />
    </React.Fragment>
  );
};

export default XFileInput;
