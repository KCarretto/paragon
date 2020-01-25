import { useQuery } from "@apollo/react-hooks";
import { ApolloError } from "apollo-client/errors/ApolloError";
import gql from "graphql-tag";
import moment from "moment";
import * as React from "react";
import { useState } from "react";
import { useParams } from "react-router-dom";
import { Card, Container, Label } from "semantic-ui-react";
import { XCredentialSummary } from "../components/credential";
import { XErrorMessage, XLoadingMessage } from "../components/messages";
import { XTargetHeader } from "../components/target";
import { XTaskSummary } from "../components/task";
import { Credential, Tag, Target, Task } from "../graphql/models";

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

  const [name, setName] = useState<string>(null);
  const [primaryIP, setPrimaryIP] = useState<string>(null);
  const [publicIP, setPublicIP] = useState<string>(null);
  const [machineUUID, setMachineUUID] = useState<string>(null);
  const [primaryMAC, setPrimaryMAC] = useState<string>(null);
  const [hostname, setHostname] = useState<string>(null);
  const [lastSeen, setLastSeen] = useState<any>(null);
  const [tasks, setTasks] = useState<Task[]>([]);
  const [tags, setTags] = useState<Tag[]>([]);
  const [creds, setCreds] = useState<Credential[]>([]);

  const { called, loading } = useQuery(TARGET_QUERY, {
    variables: { id },
    pollInterval: 5000,
    onCompleted: (data: TargetQueryResponse) => {
      setError(null);

      setName(data.target.name);
      setPrimaryIP(data.target.primaryIP);
      setPublicIP(data.target.publicIP);
      setPrimaryMAC(data.target.primaryMAC);
      setMachineUUID(data.target.machineUUID);
      setHostname(data.target.hostname);
      setLastSeen(data.target.lastSeen);
      setTasks(data.target.tasks || []);
      setTags(data.target.tags || []);
      setCreds(data.target.credentials || []);
    },
    onError: err => setError(err)
  });

  return (
    <Container fluid style={{ padding: "20px" }}>
      <XTargetHeader name={name} tags={tags} />

      <XErrorMessage title="Error Loading Target" err={error} />
      <XLoadingMessage
        title="Loading Target"
        msg="Fetching target information..."
        hidden={called && !loading}
      />
      <Card fluid centered>
        <Card.Content>
          {lastSeen &&
          moment(lastSeen).isBefore(moment().subtract(5, "minutes")) ? (
            <Label
              corner="right"
              size="large"
              icon="times circle"
              color="red"
            />
          ) : (
            <Label
              corner="right"
              size="large"
              icon="check circle"
              color="green"
            />
          )}
          <Card.Meta>
            <a>
              <i aria-hidden="true" className="clock icon"></i>
              Last Seen: {lastSeen ? moment(lastSeen).fromNow() : "Never"}
              <br />
            </a>
            {primaryIP ? (
              <a>
                <i aria-hidden="true" className="user icon"></i>
                Primary IP: {primaryIP}
                <br />
              </a>
            ) : (
              <div></div>
            )}
            {hostname ? (
              <a>
                <i aria-hidden="true" className="user icon"></i>
                Hostname: {hostname}
                <br />
              </a>
            ) : (
              <div></div>
            )}
          </Card.Meta>
          <Card.Description>
            <XTaskSummary tasks={tasks} limit={tasks.length} />
            <XCredentialSummary {...creds} />
          </Card.Description>
        </Card.Content>
        {primaryMAC || publicIP || machineUUID ? (
          <Card.Content extra>
            {primaryMAC ? (
              <a>
                <i aria-hidden="true" className="user icon"></i>
                Primary MAC: {primaryMAC}
                <br />
              </a>
            ) : (
              <div></div>
            )}
            {publicIP ? (
              <a>
                <i aria-hidden="true" className="user icon"></i>
                Public IP: {publicIP}
                <br />
              </a>
            ) : (
              <div></div>
            )}
            {machineUUID ? (
              <a>
                <i aria-hidden="true" className="user icon"></i>
                MachineUUID: {machineUUID}
                <br />
              </a>
            ) : (
              <div></div>
            )}
          </Card.Content>
        ) : (
          <div></div>
        )}
      </Card>
    </Container>
  );
};

export default XTargetView;
