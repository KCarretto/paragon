import { useMutation, useQuery } from "@apollo/react-hooks";
import gql from "graphql-tag";
import * as React from "react";
import { FunctionComponent } from "react";
import { Link } from "react-router-dom";
import { Button, Label, Table } from "semantic-ui-react";
import { XNoCredentialsFound } from "../components/credential";
import { XBoundary } from "../components/layout";
import { XErrorMessage, XLoadingMessage } from "../components/messages";
import { Credential } from "../graphql/models";

export const MULTI_CREDENTIAL_QUERY = gql`
  {
    credentials {
      id
      principal
      secret
      kind
      fails
      target {
        id
        name
      }
    }
  }
`;

const REMOVE_CREDENTIAL = gql`
  mutation RemoveCredential($credential: ID!) {
    deleteCredential(input: { id: $credential })
  }
`;

const TargetLabel: FunctionComponent<{ credential: Credential }> = ({
  credential
}) => {
  const [removeCredential] = useMutation<any, { credential: string }>(
    REMOVE_CREDENTIAL,
    {
      refetchQueries: [{ query: MULTI_CREDENTIAL_QUERY }]
    }
  );

  const handleRemove = () => {
    let vars = {
      credential: credential.id
    };

    removeCredential({
      variables: vars
    });
  };

  return (
    <Button.Group labeled>
      <Button compact as={Link} to={"/targets/" + credential.target.id}>
        {`${credential.target.name} (fails: ${credential.fails})`}
      </Button>
      <Button compact negative icon="x" onClick={handleRemove} />
    </Button.Group>
  );
};

export type MultiCredentialResponse = {
  credentials: Credential[];
};

type XCredentialTableRowProps = {
  credentials: [Credential];
};

const XCredentialTableRow: FunctionComponent<XCredentialTableRowProps> = ({
  credentials
}) => {
  const whenEmpty = <span>No Credentials</span>;
  const credential = credentials[0];
  return (
    <Table.Row>
      <Table.Cell width={2}>
        <b>{credential.principal}</b>
      </Table.Cell>
      <Table.Cell collapsing width={8}>
        <b>{credential.secret}</b>
      </Table.Cell>
      <Table.Cell collapsing width={8}>
        {credential.kind}
      </Table.Cell>
      <Table.Cell collapsing width={8}>
        <Label.Group style={{ maxWidth: "55vw", overflowX: "auto" }}>
          {credentials.map((credential, index) => (
            <TargetLabel key={index} credential={credential} />
          ))}
        </Label.Group>
      </Table.Cell>
    </Table.Row>
  );
};

const XMultiCredentialView = () => {
  const { loading, error, data: { credentials = [] } = {} } = useQuery<
    MultiCredentialResponse
  >(MULTI_CREDENTIAL_QUERY, {
    pollInterval: 5000
  });

  const whenLoading = (
    <XLoadingMessage
      title="Loading Credentials"
      msg="Fetching Credential info"
    />
  );
  const whenEmpty = (
    <Table.Row>
      <Table.HeaderCell colSpan="4">
        <XNoCredentialsFound />
      </Table.HeaderCell>
    </Table.Row>
  );

  const aggregatedCredentials = {};
  credentials &&
    credentials.forEach(credential => {
      const key = credential.principal + credential.secret;
      aggregatedCredentials[key]
        ? aggregatedCredentials[key].push(credential)
        : (aggregatedCredentials[key] = [credential]);
    });
  return (
    <React.Fragment>
      <XErrorMessage title="Error Loading Credentials" err={error} />

      <XBoundary boundary={whenLoading} show={!loading}>
        <Table celled style={{ overflow: "auto" }}>
          <Table.Header>
            <Table.Row>
              <Table.HeaderCell>Credential</Table.HeaderCell>
              <Table.HeaderCell>Secret</Table.HeaderCell>
              <Table.HeaderCell>Type</Table.HeaderCell>
              <Table.HeaderCell>Targets</Table.HeaderCell>
            </Table.Row>
          </Table.Header>

          <Table.Body>
            <XBoundary
              boundary={whenEmpty}
              show={credentials && credentials.length > 0}
            >
              {credentials &&
                Object.keys(aggregatedCredentials).map((key, index) => (
                  <XCredentialTableRow
                    key={index}
                    credentials={aggregatedCredentials[key]}
                  />
                ))}
            </XBoundary>
          </Table.Body>

          <Table.Footer fullWidth></Table.Footer>
        </Table>
      </XBoundary>
    </React.Fragment>
  );
};

export default XMultiCredentialView;
