import { ApolloError } from "apollo-client/errors/ApolloError";
import * as React from "react";
import { useState } from "react";
import { Button, ButtonProps, Form, Grid, Input, Modal } from "semantic-ui-react";
import { HTTP_URL } from "../../config";
import { useModal, XFileInput } from "../form";
import { XErrorMessage } from "../messages";

type FileUploadModalParams = {
  fileName?: string;
  button?: ButtonProps;
  openOnStart?: boolean;
};

const XFileUploadModal = ({
  openOnStart,
  button,
  fileName
}: FileUploadModalParams) => {
  const [openModal, closeModal, isOpen] = useModal();
  const [error, setError] = useState<ApolloError | null>(null);

  // Form params
  const [name, setName] = useState<string>(fileName || "");
  const [content, setContent] = useState<File>(null);
  const [loading, setLoading] = useState<boolean>(false);

  const handleSubmit = () => {
    var data = new FormData();
    data.append("fileName", name);
    data.append("fileContent", content);
    setLoading(true);

    fetch(HTTP_URL + "/cdn/upload/", {
      mode: "no-cors",
      method: "POST",
      body: data
    })
      .then(resp => {
        setLoading(false);
        setError(null);
        setName(null);
        setContent(null);
        closeModal();
      })
      .catch(err => {
        let e = new ApolloError({ errorMessage: String(err) });
        setLoading(false);
        setError(e);
      });
  };

  if (openOnStart) {
    openModal();
  }

  return (
    <Modal
      open={isOpen}
      onClose={closeModal}
      trigger={<Button onClick={openModal} {...button} />}
      // trigger={React.createElement(trigger, {onClick={openModal}}) || <Button positive circular icon="plus" onClick={openModal} />}
      size="large"
      // Form properties
      as={Form}
      onSubmit={handleSubmit}
    >
      <Modal.Header>Upload a File</Modal.Header>
      <Modal.Content>
        <Grid verticalAlign="middle" stackable container columns={"equal"}>
          <Grid.Column>
            <Input
              label="File Name"
              icon="file"
              fluid
              placeholder="Enter file name"
              name="fileName"
              value={name}
              onChange={(e, { value }) => setName(value)}
            />
          </Grid.Column>
          <Grid.Column>
            <XFileInput
              id="upload_fileContent"
              setFile={f => setContent(f)}
              file={content}
            />
          </Grid.Column>
        </Grid>
        <XErrorMessage title="Failed to upload file" err={error} />
      </Modal.Content>
      <Modal.Actions>
        <Form.Button loading={loading} style={{ marginBottom: "10px" }} positive floated="right">
          Upload
        </Form.Button>
      </Modal.Actions>
    </Modal>
  );
};

export default XFileUploadModal;
