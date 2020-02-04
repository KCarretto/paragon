import { useQuery } from "@apollo/react-hooks";
import { ApolloError } from "apollo-client/errors/ApolloError";
import gql from "graphql-tag";
import * as React from "react";
import { useState } from "react";
import { Link, useParams } from "react-router-dom";
import { Button, Container, Icon, Label } from "semantic-ui-react";
import { XJobHeader } from "../components/job";
import { XErrorMessage, XLoadingMessage } from "../components/messages";
import {
  XTaskContent,
  XTaskError,
  XTaskOutput,
  XTaskStatus
} from "../components/task";
import { Tag, Target, Task } from "../graphql/models";

const TASK_QUERY = gql`
  query Task($id: ID!) {
    task(id: $id) {
      id
      queueTime
      claimTime
      execStartTime
      execStopTime
      content
      output
      error
      sessionID

      target {
        id
        name
      }

      job {
        id
        name
        tags {
          id
          name
        }
      }
    }
  }
`;

type TaskQueryResponse = {
  task: Task;
};

const XTaskView = () => {
  let { id } = useParams();
  const [loadingError, setLoadingError] = useState<ApolloError>(null);

  const [queueTime, setQueueTime] = useState<string>(null);
  const [claimTime, setClaimTime] = useState<string>(null);
  const [execStartTime, setExecStartTime] = useState<string>(null);
  const [execStopTime, setExecStopTime] = useState<string>(null);
  const [content, setContent] = useState<string>(null);
  const [output, setOutput] = useState<string>(null);
  const [error, setError] = useState<string>(null);
  const [sessionID, setSessionID] = useState<string>(null);
  const [target, setTarget] = useState<Target>({});
  const [jobID, setJobID] = useState<string>(null);
  const [name, setName] = useState<string>("");
  const [tags, setTags] = useState<Tag[]>([]);

  const { called, loading } = useQuery<TaskQueryResponse>(TASK_QUERY, {
    variables: { id },
    onCompleted: data => {
      setLoadingError(null);

      if (!data || !data.task) {
        data = {
          task: {
            id: null,
            queueTime: null,
            claimTime: null,
            execStartTime: null,
            execStopTime: null,
            content: null,
            output: null,
            error: null,
            sessionID: null
          }
        };
      }
      if (!data.task.target) {
        data.task.target = {};
      }
      if (!data.task.job) {
        data.task.job = { id: null, name: "", tags: [] };
      }

      setQueueTime(data.task.queueTime);
      setClaimTime(data.task.claimTime);
      setExecStartTime(data.task.execStartTime);
      setExecStopTime(data.task.execStopTime);
      setContent(data.task.content);
      setOutput(data.task.output);
      setError(data.task.error);
      setSessionID(data.task.sessionID);
      setTarget(data.task.target);
      setJobID(data.task.job.id);
      setTags(data.task.job.tags);
      setName(data.task.job.name);
    },
    onError: err => setLoadingError(err)
  });

  let status = new XTaskStatus().getStatus({
    id: id,
    queueTime: queueTime,
    claimTime: claimTime,
    execStartTime: execStartTime,
    execStopTime: execStopTime,
    error: error
  }).icon;

  return (
    <Container fluid style={{ padding: "20px" }}>
      <Link to={"/jobs/" + jobID}>
        <XJobHeader
          name={name}
          tags={tags}
          icon={React.createElement(Icon, { size: "large", ...status })}
        />
      </Link>
      {!target || !target.id ? (
        <span />
      ) : (
        <Button
          basic
          animated
          color="blue"
          size="small"
          style={{ margin: "15px" }}
          as={Link}
          to={"/targets/" + target.id}
        >
          <Button.Content visible>
            {target.name || "View Target"}
          </Button.Content>
          <Button.Content hidden>
            <Icon name="arrow right" />
          </Button.Content>
        </Button>
      )}
      {!sessionID ? <span /> : <Label>SessionID: {sessionID}</Label>}

      <XErrorMessage title="Error Loading Task" err={loadingError} />
      <XLoadingMessage
        title="Loading Task"
        msg="Fetching task information..."
        hidden={called && !loading}
      />

      <XTaskContent content={content} />
      <XTaskOutput output={output} />
      <XTaskError error={error} />
    </Container>
  );
};

export default XTaskView;
