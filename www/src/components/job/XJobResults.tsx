import { useQuery } from "@apollo/react-hooks";
import gql from "graphql-tag";
import * as React from "react";
import { useState } from "react";
import { Job, Task } from "../../graphql/models";
import { XBoundary, XCardGroup } from "../layout";
import { XErrorMessage, XLoadingMessage } from "../messages";
import { XTaskResultCard } from "../task";

export const RUNS_QUERY = gql`
  query Runs($name: String!) {
    jobs(input: { search: $name }) {
      id

      tasks {
        id
        queueTime
        claimTime
        execStartTime
        execStopTime
        error

        target {
          id
          name
        }
      }
    }
  }
`;

const XJobResults: React.FC<{ name: string }> = ({ name }) => {
  const {
    loading,
    error,
    data: { jobs = [] } = {
      jobs: []
    }
  } = useQuery<{ jobs: Job[] }, { name: string }>(RUNS_QUERY, {
    variables: { name: name },
    pollInterval: 2500
  });

  // Map of Target Name => []Task
  const taskMap = new Map<string, Task[]>();
  jobs.forEach(job => {
    job.tasks.forEach(task => {
      let key = task.target.name;
      if (!taskMap.has(key)) {
        taskMap[key] = [];
      }
      taskMap[key].push(task);
    });
  });

  const [displayIndex, setDisplayIndex] = useState<string>(null);

  return (
    <React.Fragment>
      <XBoundary
        boundary={
          <XLoadingMessage
            title="Fetching Results"
            msg="Loading execution results..."
          />
        }
        show={!loading}
      >
        <XBoundary
          boundary={
            <XErrorMessage title="Failed to Load Results" err={error} />
          }
          show={!error}
        >
          {!displayIndex ? (
            <XCardGroup>
              {Array.from(taskMap).map(([targetName, tasks], index) => (
                <XTaskResultCard
                  key={index}
                  targetName={targetName}
                  tasks={tasks}
                />
              ))}
            </XCardGroup>
          ) : (
            <span id="OUTPUT" />
          )}
        </XBoundary>
      </XBoundary>
    </React.Fragment>
  );
};

export default XJobResults;
