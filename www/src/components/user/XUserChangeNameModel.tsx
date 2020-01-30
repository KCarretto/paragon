import { useMutation } from "@apollo/react-hooks";
import { ApolloError } from "apollo-client/errors/ApolloError";
import gql from "graphql-tag";
import * as React from "react";
import { useState } from "react";
import { Button, Form, Grid, Input, Message, Modal } from "semantic-ui-react";
import { PROFILE_QUERY } from "../../views/XProfileView";
import { useModal } from "../form";

export const CHANGE_NAME_MUTATION = gql`
  mutation changeName($name: String!) {
    changeName(input: { name: $name }) {
      id
      name
      photoURL
      isActivated
      isAdmin
    }
  }
`;

type ChangeNameModalParams = {
  openOnStart?: boolean;
};

const XUserChangeNameModal = ({ openOnStart }: ChangeNameModalParams) => {
  const [openModal, closeModal, isOpen] = useModal();
  const [error, setError] = useState<ApolloError>(null);

  // Form params
  const [name, setName] = useState<string>("");

  const [changeName, { called, loading }] = useMutation(CHANGE_NAME_MUTATION, {
    refetchQueries: [{ query: PROFILE_QUERY }]
  });

  const handleSubmit = () => {
    let vars = {
      name: name
    };

    changeName({
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
      .catch(err => setError(err));
  };

  if (openOnStart) {
    openModal();
  }

  return (
    <Modal
      open={isOpen}
      onClose={closeModal}
      trigger={
        <Button positive fluid onClick={openModal}>
          Change Name
        </Button>
      }
      size="large"
      // Form properties
      as={Form}
      onSubmit={handleSubmit}
      error={called && error}
      loading={called && loading}
    >
      <Modal.Header>{"Change your Name"}</Modal.Header>
      <Modal.Content>
        <Grid verticalAlign="middle" stackable container columns={"equal"}>
          <Grid.Column>
            <Input
              fluid
              placeholder="Alex Smith"
              name="name"
              value={name}
              onChange={(e, { value }) => setName(value)}
            />
          </Grid.Column>
        </Grid>

        <Message
          error
          icon="warning"
          header={"Failed to Change Name"}
          onDismiss={(e, data) => setError(null)}
          content={error ? error.message : "Unknown Error"}
        />
      </Modal.Content>
      <Modal.Actions>
        <Form.Button style={{ marginBottom: "10px" }} positive floated="right">
          Send It!
        </Form.Button>
      </Modal.Actions>
    </Modal>
  );
};

export default XUserChangeNameModal;
