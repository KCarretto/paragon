import { useQuery } from "@apollo/react-hooks";
import gql from "graphql-tag";
import moment from "moment";
import * as React from "react";
import { useParams } from "react-router-dom";
import { Header, Icon, Table } from "semantic-ui-react";
import { XCredentialSummary } from "../components/credential";
import { XClipboard } from "../components/form";
import { XBoundary, XCardGroup } from "../components/layout";
import { XErrorMessage, XLoadingMessage } from "../components/messages";
import { XTargetHeader } from "../components/target";
import {
  XNoTasksFound,
  XTaskCard,
  XTaskCardDisplayType
} from "../components/task";
import { Target } from "../graphql/models";

export const TARGET_QUERY = gql`
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

  const {
    loading,
    error,
    data: {
      target: {
        name = "Untitled Target",
        primaryIP = null,
        primaryMAC = null,
        hostname = null,
        machineUUID = null,
        lastSeen = null,
        tags = [],
        tasks = [],
        credentials = []
      } = {}
    } = {}
  } = useQuery<TargetQueryResponse>(TARGET_QUERY, {
    variables: { id },
    pollInterval: 5000
  });

  const whenLoading = (
    <XLoadingMessage title="Loading Target" msg="Fetching target info" />
  );
  const whenFieldEmpty = <span>Unknown</span>;
  const whenNotSeen = <span>Never</span>;
  const whenTasksEmpty = <XNoTasksFound />;

  return (
    <React.Fragment>
      <XTargetHeader name={name} tags={tags} lastSeen={lastSeen} />

      <XErrorMessage title="Error Loading Target" err={error} />
      <XBoundary boundary={whenLoading} show={!loading}>
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
                <XBoundary boundary={whenFieldEmpty} show={hostname !== ""}>
                  <a>
                    <XClipboard value={hostname}>{hostname}</XClipboard>
                  </a>
                </XBoundary>
              </Table.Cell>
            </Table.Row>
            <Table.Row>
              <Table.HeaderCell collapsing>
                <Icon name="time" style={{ marginLeft: "10px" }} />
                Last Seen
              </Table.HeaderCell>
              <Table.Cell>
                <XBoundary boundary={whenNotSeen} show={lastSeen}>
                  <a>
                    <XClipboard value={lastSeen}>
                      {" "}
                      {moment(lastSeen).fromNow()}
                    </XClipboard>
                  </a>
                </XBoundary>
              </Table.Cell>
            </Table.Row>
            <Table.Row>
              <Table.HeaderCell collapsing>
                <Icon name="wifi" style={{ marginLeft: "10px" }} />
                Primary IP
              </Table.HeaderCell>
              <Table.Cell>
                <XBoundary boundary={whenFieldEmpty} show={primaryIP !== ""}>
                  <a>
                    <XClipboard value={primaryIP}>{primaryIP}</XClipboard>
                  </a>
                </XBoundary>
              </Table.Cell>
            </Table.Row>
            <Table.Row>
              <Table.HeaderCell collapsing>
                <Icon name="microchip" style={{ marginLeft: "10px" }} />
                Primary MAC
              </Table.HeaderCell>
              <Table.Cell>
                <XBoundary boundary={whenFieldEmpty} show={primaryMAC !== ""}>
                  <a>
                    <XClipboard value={primaryMAC}>{primaryMAC}</XClipboard>
                  </a>
                </XBoundary>
              </Table.Cell>
            </Table.Row>
            <Table.Row>
              <Table.HeaderCell collapsing>
                <Icon name="id card outline" style={{ marginLeft: "10px" }} />
                MachineUUID
              </Table.HeaderCell>
              <Table.Cell>
                <XBoundary boundary={whenFieldEmpty} show={machineUUID !== ""}>
                  <a>
                    <XClipboard value={machineUUID}>{machineUUID}</XClipboard>
                  </a>
                </XBoundary>
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

        <XBoundary boundary={whenTasksEmpty} show={tasks.length > 0}>
          <XCardGroup>
            {tasks.map(task => (
              <XTaskCard
                key={task.id}
                display={XTaskCardDisplayType.JOB}
                task={task}
              />
            ))}
          </XCardGroup>
        </XBoundary>
      </XBoundary>
    </React.Fragment>
  );
};

export default XTargetView;
