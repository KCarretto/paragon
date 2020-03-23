import { useQuery } from "@apollo/react-hooks";
import gql from "graphql-tag";
import * as React from "react";
import { Link, useParams } from "react-router-dom";
import { Button, Icon, Label } from "semantic-ui-react";
import { XJobHeader } from "../components/job";
import { XBoundary } from "../components/layout";
import { XErrorMessage, XLoadingMessage } from "../components/messages";
import {
  XTaskContent,
  XTaskError,
  XTaskOutput,
  XTaskStatus
} from "../components/task";
import { Task } from "../graphql/models";

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
        staged

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

  const {
    loading,
    error,
    data: {
      task: {
        queueTime = null,
        claimTime = null,
        execStartTime = null,
        execStopTime = null,
        content = null,
        output = null,
        error: outputErr = null,
        sessionID = "",
        target: targetData = { id: null, name: null },
        job: jobData = { id: null, name: "", tags: [] }
      } = {}
    } = {}
  } = useQuery<TaskQueryResponse>(TASK_QUERY, {
    variables: { id }
  });
  const target = targetData || { id: null, name: null };
  const job = jobData || { id: null, name: "", tags: [] };

  let status = new XTaskStatus().getStatus({
    id: id,
    queueTime: queueTime,
    claimTime: claimTime,
    execStartTime: execStartTime,
    execStopTime: execStopTime,
    error: outputErr
  }).icon;

  const whenLoading = (
    <XLoadingMessage title="Loading Task" msg="Fetching task info" />
  );

  return (
    <React.Fragment>
      <XErrorMessage title="Error Loading Task" err={error} />
      <XBoundary boundary={whenLoading} show={!loading}>
        <Link to={"/jobs/" + job.id || "0"}>
          <XJobHeader
            name={job.name}
            tags={job.tags}
            icon={React.createElement(Icon, { size: "large", ...status })}
          />
        </Link>

        <XBoundary boundary={<span />} show={target.id !== null}>
          <Button
            basic
            animated
            color="blue"
            size="small"
            style={{ margin: "15px" }}
            as={Link}
            to={"/targets/" + target.id || "0"}
          >
            <Button.Content visible>
              {target.name || "View Target"}
            </Button.Content>
            <Button.Content hidden>
              <Icon name="arrow right" />
            </Button.Content>
          </Button>
        </XBoundary>

        <XBoundary boundary={<span />} show={sessionID !== ""}>
          <Label>SessionID: {sessionID}</Label>
        </XBoundary>

        <XTaskContent content={content} />
        <XTaskOutput output={output} />
        <XTaskError error={outputErr} />
      </XBoundary>
    </React.Fragment>
  );
};

export default XTaskView;
