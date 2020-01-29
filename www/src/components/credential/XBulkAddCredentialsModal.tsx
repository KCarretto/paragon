import { useMutation } from "@apollo/react-hooks";
import gql from "graphql-tag";
import * as React from "react";
import { useState } from "react";
import { Button, Form, Grid, Input, Modal } from "semantic-ui-react";
import { Target } from "../../graphql/models";
import { MULTI_TARGET_QUERY } from "../../views";
import { useModal, XTargetTypeahead } from "../form";
import { XErrorMessage } from "../messages";

export const BULK_ADD_CREDS_MUTATION = gql`
  mutation AddCredentialForTargets(
    $principal: String!
    $secret: String!
    $targets: [ID!]
    $kind: String
  ) {
    addCredentialForTargets(
      input: {
        ids: $targets
        principal: $principal
        secret: $secret
        kind: $kind
      }
    ) {
      id
      credentials {
        id
        principal
        secret
      }
    }
  }
`;

const XBulkAddCredentialsModal = () => {
  const [openModal, closeModal, isOpen] = useModal();
  // const [error, setError] = useState<ApolloError>(null);

  // Form params
  const [principal, setPrincipal] = useState<string>("");
  const [secret, setSecret] = useState<string>("");
  const [targets, setTargets] = useState<Target[]>([]);

  const [addCredentials, { called, loading, error }] = useMutation(
    BULK_ADD_CREDS_MUTATION,
    {
      refetchQueries: [{ query: MULTI_TARGET_QUERY }]
    }
  );

  const handleSubmit = () => {
    let vars = {
      principal: principal,
      secret: secret,
      targets: targets
    };

    addCredentials({
      variables: vars
    }).then(closeModal);
  };

  return (
    <Modal
      open={isOpen}
      onClose={closeModal}
      trigger={<Button positive circular icon="key" onClick={openModal} />}
      size="large"
      // Form properties
      as={Form}
      onSubmit={handleSubmit}
      error={called && error}
      loading={called && loading}
    >
      <Modal.Header>{"Add Credentials for Targets"}</Modal.Header>
      <Modal.Content>
        <Grid verticalAlign="middle" stackable container columns={"equal"}>
          <Grid.Column>
            <Input
              label="Principal"
              icon="user"
              fluid
              placeholder="Enter principal (i.e. username)"
              name="principal"
              value={principal}
              onChange={(e, { value }) => setPrincipal(value)}
            />
          </Grid.Column>
          <Grid.Column>
            <Input
              label="Secret"
              icon="key"
              fluid
              placeholder="Enter secret (i.e. password)"
              name="secret"
              value={secret}
              onChange={(e, { value }) => setSecret(value)}
            />
          </Grid.Column>
          <Grid.Column>
            <XTargetTypeahead
              labeled
              onChange={(e, { value }) => setTargets(value)}
            />
          </Grid.Column>
        </Grid>

        <XErrorMessage title="Failed to Add Credentials" err={error} />
      </Modal.Content>
      <Modal.Actions>
        <Form.Button style={{ marginBottom: "10px" }} positive floated="right">
          Add
        </Form.Button>
      </Modal.Actions>
    </Modal>
  );
};

export default XBulkAddCredentialsModal;
