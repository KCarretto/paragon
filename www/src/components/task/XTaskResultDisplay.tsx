import { useQuery } from "@apollo/react-hooks";
import gql from "graphql-tag";
import moment from "moment";
import * as React from "react";
import { Link } from "react-router-dom";
import { Button, Header, Icon, Label, Segment } from "semantic-ui-react";
import { XTaskStatus } from ".";
import { Task } from "../../graphql/models";
import { XBoundary } from "../layout";
import { XErrorMessage, XLoadingMessage } from "../messages";

const RESULT_QUERY = gql`
  query Result($id: ID!) {
    task(id: $id) {
      id
      output
      error
      queueTime
      claimTime
      execStartTime
      execStopTime
      error

      target {
        id
        name
      }

      job {
        id
        staged
      }
    }
  }
`;

const XTaskResultDisplay: React.FC<{
  id: string;
  version: string;
  onExit: () => void;
}> = ({ id, version, onExit }) => {
  const {
    loading,
    error,
    data: { task = null } = {
      task: null
    }
  } = useQuery<{ task: Task }, { id: string }>(RESULT_QUERY, {
    skip: !id,
    variables: { id: id },
    pollInterval: 2500
  });

  return (
    <XBoundary
      boundary={
        <XLoadingMessage
          title="Loading Result"
          msg="Fetching task execution output"
        />
      }
      show={!loading}
    >
      <XBoundary
        boundary={<XErrorMessage title="Failed to Load Result" err={error} />}
        show={!error}
      >
        {task && (
          <React.Fragment>
            <Header>
              <Header.Content>
                <Button icon="arrow left" color="teal" onClick={onExit} />
                <Icon
                  floated="right"
                  size="large"
                  {...new XTaskStatus().getStatus(task).icon}
                />
                {task.target && <Link to={'/targets/' + task.target.id}>{task.target.name}</Link>}
                {" "}<Label>{version}</Label>
              </Header.Content>
              <Header.Subheader style={{ paddingLeft: "3em", paddingTop: "1em" }}>
                Last Updated:{" "}{moment(new XTaskStatus().getTimestamp(task)).fromNow()}
              </Header.Subheader>
            </Header>
            <Segment.Group raised padded style={{ overflow: "auto" }}>
              {task.error && <Segment color="red" >
                <Header color="red" content="Error" />
                <pre style={{ color: "red" }}>{task.error}</pre>
              </Segment>}
              {task.output && <Segment color="blue">
                <Header content="Output" />
                <pre>{task.output}</pre>
              </Segment>}
            </Segment.Group>
          </React.Fragment>
        )}
      </XBoundary>
    </XBoundary>
  );
};
export default XTaskResultDisplay;
