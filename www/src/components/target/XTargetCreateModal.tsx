import { useMutation } from "@apollo/react-hooks";
import { ApolloError } from "apollo-client/errors/ApolloError";
import gql from "graphql-tag";
import * as React from "react";
import { useState } from "react";
import { Button, Form, Grid, Input, Message, Modal } from "semantic-ui-react";
import { Tag } from "../../graphql/models";
import { MULTI_TARGET_QUERY } from "../../views";
import { useModal, XTagTypeahead } from "../form";

export const CREATE_TARGET_MUTATION = gql`
  mutation CreateTarget($name: String!, $primaryIP: String!, $tags: [ID!]) {
    createTarget(input: { name: $name, primaryIP: $primaryIP, tags: $tags }) {
      id
      name
      primaryIP
      lastSeen
      tags {
        id
        name
      }
    }
  }
`;

const XTargetCreateModal = () => {
  const [openModal, closeModal, isOpen] = useModal();
  const [error, setError] = useState<ApolloError>(null);

  // Form Params
  const [name, setName] = useState<string>("");
  const [primaryIP, setPrimaryIP] = useState<string>("");
  const [tags, setTags] = useState<Tag[]>([]);

  const [createTarget, { called, loading }] = useMutation(
    CREATE_TARGET_MUTATION,
    {
      refetchQueries: [{ query: MULTI_TARGET_QUERY }]
    }
  );

  const handleSubmit = () => {
    let vars = {
      name: name,
      primaryIP: primaryIP,
      tags: tags
    };

    createTarget({
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
      size="small"
      // Form properties
      as={Form}
      onSubmit={handleSubmit}
      error={called && error}
      loading={called && loading}
    >
      <Modal.Header>Create a Target</Modal.Header>
      <Modal.Content>
        <Grid verticalAlign="middle" centered stackable container columns={2}>
          <Grid.Column>
            <Input
              label="Target Name"
              icon="desktop"
              fluid
              placeholder="Enter target name"
              name="name"
              value={name}
              onChange={(e, { value }) => setName(value)}
            />
          </Grid.Column>

          <Grid.Column>
            <Input
              label="Primary IP"
              icon="desktop"
              fluid
              placeholder="Enter primary ip address"
              name="primaryIP"
              value={primaryIP}
              onChange={(e, { value }) => setPrimaryIP(value)}
            />
          </Grid.Column>

          <Grid.Column>
            <XTagTypeahead
              labeled
              onChange={(e, { value }) => setTags(value)}
            />
          </Grid.Column>
        </Grid>
        <Message
          error
          icon="warning"
          header={"Failed to Create Target"}
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

export default XTargetCreateModal;
