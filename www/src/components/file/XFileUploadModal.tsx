import { ApolloError } from "apollo-client/errors/ApolloError";
import * as React from "react";
import { useState } from "react";
import { Button, Form, Grid, Input, Message, Modal } from "semantic-ui-react";
import { InputFile } from "semantic-ui-react-input-file";
import { HTTP_URL } from "../../config";
import { useModal } from "../form";

type FileUploadModalParams = {
  openOnStart?: boolean;
};

const XFileUploadModal = ({ openOnStart }: FileUploadModalParams) => {
  const [openModal, closeModal, isOpen] = useModal();
  const [error, setError] = useState<ApolloError>(null);

  // Form params
  const [name, setName] = useState<string>("");
  const [content, setContent] = useState<string>(null);

  const handleSubmit = () => {
    var data = new FormData();
    data.append("fileName", name);
    data.append("fileContent", content);

    fetch(HTTP_URL + "/cdn/upload/", {
      mode: "no-cors",
      method: "POST",
      body: data
    }).then(
      function(res) {
        if (res.ok) {
          closeModal();
          alert("Perfect! ");
        } else if (res.status === 401) {
          alert("Oops! ");
        }
      },
      function(e) {
        alert("Error submitting form!");
      }
    );
  };

  if (openOnStart) {
    openModal();
  }

  return (
    <Modal
      open={isOpen}
      onClose={closeModal}
      trigger={<Button positive circular icon="plus" onClick={openModal} />}
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
            <InputFile
              input={{
                id: "upload_fileContent",
                onChange: e => {
                  e.preventDefault();
                  console.log("FILE UPLOAD EVENT", e);
                  console.log("FILE UPLOAD EVENT TARGET", e.target);
                  console.log("FILE UPLOAD EVENT TARGET VALUE", e.target.value);
                  console.log("FILE UPLOAD EVENT TARGET FILES", e.target.files);
                  setContent(e.target.files[0]);
                }
              }}
            />
          </Grid.Column>
        </Grid>

        <Message
          error
          icon="warning"
          header={"Failed to Upload File"}
          onDismiss={(e, data) => setError(null)}
          content={error ? error.message : "Unknown Error"}
        />
      </Modal.Content>
      <Modal.Actions>
        <Form.Button style={{ marginBottom: "10px" }} positive floated="right">
          Queue
        </Form.Button>
      </Modal.Actions>
    </Modal>
  );
};

export default XFileUploadModal;
