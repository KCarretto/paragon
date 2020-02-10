import { useQuery } from "@apollo/react-hooks";
import { ApolloError } from "apollo-client/errors/ApolloError";
import gql from "graphql-tag";
import moment from "moment";
import * as React from "react";
import { useState } from "react";
import { useParams } from "react-router-dom";
import { Header, Icon, Table } from "semantic-ui-react";
import { XCredentialSummary } from "../components/credential";
import { XClipboard } from "../components/form";
import { XCardGroup } from "../components/layout";
import { XErrorMessage, XLoadingMessage } from "../components/messages";
import { XTargetHeader } from "../components/target";
import { XTaskCard, XTaskCardDisplayType } from "../components/task";
import { Target } from "../graphql/models";

const TARGET_QUERY = gql`
  query Target($id: ID!) {
    target(id: $id) {
      id
      name
      primaryIP
      publicIP
      primaryMAC
      machineUUID
      hostname
      lastSeen
      tasks {
        id
        queueTime
        claimTime
        execStartTime
        execStopTime
        error
        job {
          id
          name
        }
      }
      tags {
        id
        name
      }
      credentials {
        id
        principal
        secret
        fails
      }
    }
  }
`;

type TargetQueryResponse = {
  target: Target;
};

const XTargetView = () => {
  let { id } = useParams();
  const [error, setError] = useState<ApolloError>(null);
  const [
    {
      name,
      primaryIP,
      publicIP,
      machineUUID,
      primaryMAC,
      hostname,
      lastSeen,
      tasks,
      tags,
      credentials
    },
    setTarget
  ] = useState<Target>({
    name: "Untitled Target",
    tags: [],
    tasks: [],
    credentials: []
  });

  const { called, loading } = useQuery(TARGET_QUERY, {
    variables: { id },
    pollInterval: 5000,
    onCompleted: (data: TargetQueryResponse) => {
      setError(null);
      setTarget(data.target);
    },
    onError: err => setError(err)
  });

  return (
    <React.Fragment>
      <XTargetHeader name={name} tags={tags} lastSeen={lastSeen} />

      <XErrorMessage title="Error Loading Target" err={error} />
      <XLoadingMessage
        title="Loading Target"
        msg="Fetching target information..."
        hidden={called && !loading}
      />
      <Header size="large" block inverted>
        <Header.Content>Metadata</Header.Content>
      </Header>

      <Table>
        <Table.Body>
          <Table.Row>
            <Table.HeaderCell collapsing>
              <Icon name="desktop" style={{ marginLeft: "10px" }} />
              Hostname
            </Table.HeaderCell>
            <Table.Cell>
              {hostname ? (
                <a>
                  <XClipboard value={hostname}>{hostname}</XClipboard>
                </a>
              ) : (
                "Unknown"
              )}
            </Table.Cell>
          </Table.Row>
          <Table.Row>
            <Table.HeaderCell collapsing>
              <Icon name="time" style={{ marginLeft: "10px" }} />
              Last Seen
            </Table.HeaderCell>
            <Table.Cell>
              {lastSeen ? (
                <a>
                  <XClipboard value={lastSeen}>
                    {moment(lastSeen).fromNow()}
                  </XClipboard>
                </a>
              ) : (
                "Never"
              )}
            </Table.Cell>
          </Table.Row>
          <Table.Row>
            <Table.HeaderCell collapsing>
              <Icon name="wifi" style={{ marginLeft: "10px" }} />
              Primary IP
            </Table.HeaderCell>
            <Table.Cell>
              {primaryIP ? (
                <a>
                  <XClipboard value={primaryIP}>{primaryIP}</XClipboard>
                </a>
              ) : (
                "Unknown"
              )}
            </Table.Cell>
          </Table.Row>
          <Table.Row>
            <Table.HeaderCell collapsing>
              <Icon name="microchip" style={{ marginLeft: "10px" }} />
              Primary MAC
            </Table.HeaderCell>
            <Table.Cell>
              {primaryMAC ? (
                <a>
                  <XClipboard value={primaryMAC}>{primaryMAC}</XClipboard>
                </a>
              ) : (
                "Unknown"
              )}
            </Table.Cell>
          </Table.Row>
          <Table.Row>
            <Table.HeaderCell collapsing>
              <Icon name="id card outline" style={{ marginLeft: "10px" }} />
              MachineUUID
            </Table.HeaderCell>
            <Table.Cell>
              {machineUUID ? (
                <a>
                  <XClipboard value={machineUUID}>{machineUUID}</XClipboard>
                </a>
              ) : (
                "Unknown"
              )}
            </Table.Cell>
          </Table.Row>
          <Table.Row>
            <Table.HeaderCell collapsing>
              <Icon name="key" style={{ marginLeft: "10px" }} />
              Credentials
            </Table.HeaderCell>
            <Table.Cell>
              <XCredentialSummary credentials={credentials} />
            </Table.Cell>
          </Table.Row>
        </Table.Body>
      </Table>

      <Header size="large" block inverted>
        <Icon name="tasks" />
        <Header.Content>Tasks</Header.Content>
      </Header>

      <XCardGroup>
        {tasks.map(task => (
          <XTaskCard
            key={task.id}
            display={XTaskCardDisplayType.JOB}
            task={task}
          />
        ))}
      </XCardGroup>
    </React.Fragment>
  );
};

export default XTargetView;
