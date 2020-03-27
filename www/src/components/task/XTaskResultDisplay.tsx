import { useQuery } from "@apollo/react-hooks";
import gql from "graphql-tag";
import moment from "moment";
import * as React from "react";
import { Icon } from "semantic-ui-react";
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

      job {
        id
        staged
      }
    }
  }
`;

const XTaskResultDisplay: React.FC<{
  id: string;
}> = ({ id }) => {
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
            <Icon
              floated="right"
              size="large"
              {...new XTaskStatus().getStatus(task).icon}
            />
            Last Updated:{" "}
            {moment(new XTaskStatus().getTimestamp(task)).fromNow()}
            <br />
            <pre>{task.output}</pre>
          </React.Fragment>
        )}
      </XBoundary>
    </XBoundary>
  );
};
export default XTaskResultDisplay;
