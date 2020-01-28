import { useMutation } from "@apollo/react-hooks";
import { ApolloError } from "apollo-client/errors/ApolloError";
import gql from "graphql-tag";
import * as React from "react";
import { useState } from "react";
import { Button, Form, Grid, Input, Message, Modal } from "semantic-ui-react";
import { Tag } from "../../graphql/models";
import { MULTI_TAG_QUERY } from "../../views";
import { SUGGEST_TAGS_QUERY, SUGGEST_TARGETS_QUERY, useModal } from "../form";

const CREATE_TAG = gql`
  mutation CreateTag($name: String!) {
    createTag(input: { name: $name }) {
      id
      name
    }
  }
`;

type CreateTagResult = {
  data: Tag;
};

type TagCreateModalParams = {
  openOnStart?: boolean;
};

const XTagCreateModal = ({ openOnStart }: TagCreateModalParams) => {
  const [error, setError] = useState<ApolloError>(null);
  const [openModal, closeModal, isOpen] = useModal();
  if (openOnStart) {
    openModal();
  }

  // Form params
  const [name, setName] = useState<string>("");

  const [createTag, { called, loading }] = useMutation(CREATE_TAG, {
    refetchQueries: [
      { query: MULTI_TAG_QUERY },
      { query: SUGGEST_TAGS_QUERY },
      { query: SUGGEST_TARGETS_QUERY }
    ]
  });

  const handleSubmit = () => {
    let vars = {
      name: name
    };

    createTag({
      variables: vars
    })
      .then(({ data, errors }) => {
        if (errors && errors.length > 0) {
          let s = errors.map(e => e.message);
          let e = new ApolloError({
            graphQLErrors: errors,
            errorMessage: s.join("\n")
          });
          setError(e);
          return;
        }
        closeModal();
      })
      .catch(error => setError(error));
  };

  return (
    <Modal
      open={isOpen}
      onClose={closeModal}
      trigger={<Button positive circular icon="plus" onClick={openModal} />}
      size="large"
      // Form properties
      as={Form}
      onSubmit={handleSubmit}
      error={called && error}
      loading={called && loading}
    >
      <Modal.Header>Create a Tag</Modal.Header>
      <Modal.Content>
        <Grid verticalAlign="middle" stackable container columns={"equal"}>
          <Grid.Column>
            <Input
              label="Tag Name"
              icon="tag"
              fluid
              placeholder="Enter tag name"
              name="name"
              value={name}
              onChange={(e, { value }) => setName(value)}
            />
          </Grid.Column>
        </Grid>
        <Message
          error
          icon="warning"
          header={"Failed to Create Tag"}
          onDismiss={(e, data) => setError(null)}
          content={error ? error.message : "Unknown Error"}
        />
      </Modal.Content>
      <Modal.Actions>
        <Form.Button style={{ marginBottom: "10px" }} positive floated="right">
          Create
        </Form.Button>
      </Modal.Actions>
    </Modal>
  );
};

export default XTagCreateModal;
