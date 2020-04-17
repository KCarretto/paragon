import { useQuery } from "@apollo/react-hooks";
import gql from "graphql-tag";
import * as React from "react";
import { useState } from "react";
import { Button, Card, Label } from "semantic-ui-react";
import { Job, Task } from "../../graphql/models";
import { XBoundary, XCardGroup } from "../layout";
import { XErrorMessage, XLoadingMessage } from "../messages";
import { XNoTasksFound, XTaskResultCard, XTaskResultDisplay } from "../task";

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

const XJobResults: React.FC<{ name?: string }> = ({ name = null }) => {
  const {
    loading,
    error,
    data: { jobs = [] } = {
      jobs: []
    }
  } = useQuery<{ jobs: Job[] }, { name: string }>(RUNS_QUERY, {
    skip: !name,
    variables: { name: name },
    pollInterval: 2500
  });

  const [taskDisplay, setTaskDisplay] = useState<{
    active: number;
    tasks: Task[];
  }>(null);

  const ResultCards = (jobs: Job[]) => {
    if (!jobs || jobs.length < 1) {
      return <XNoTasksFound />;
    }

    // Map of Target Name => []Task
    const taskMap = new Map<string, Task[]>();
    jobs.forEach(job => {
      job.tasks.forEach(task => {
        let key = task.target.name;
        if (!taskMap.has(key)) {
          taskMap.set(key, []);
        }
        taskMap.get(key).push(task);
      });
    });

    let cards = Array.from(taskMap);
    console.log("TASK MAP", cards);

    if (cards.length < 1) {
      return <XNoTasksFound />;
    }

    return (
      <React.Fragment>
        <XCardGroup>
          {cards.map(([targetName, tasks], index) => (
            <XTaskResultCard
              key={index}
              targetName={targetName}
              tasks={tasks}
              onShowResult={(active, tasks) =>
                setTaskDisplay({ active: active, tasks: tasks })
              }
            />
          ))}
        </XCardGroup>
      </React.Fragment>
    );
  };

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
          {taskDisplay ? <Card fluid style={{ marginTop: "0px" }}>
            <Card.Content>
              <Card.Header>
                Showing results for{" "}
                {taskDisplay.tasks[taskDisplay.active].target.name}{" "}
              </Card.Header>
              <Card.Meta>
                <Button icon="arrow left" onClick={() => setTaskDisplay(null)} />
                <Label basic color="blue">
                  {taskDisplay.active !== 0
                    ? `Version ${taskDisplay.tasks.length - taskDisplay.active}`
                    : `Latest (v${taskDisplay.tasks.length})`}
                </Label>
                {/* {taskDisplay.active !== 0
                ? `Version ${taskDisplay.tasks.length - taskDisplay.active}`
                  : `Latest (v${taskDisplay.tasks.length})`} */}
              </Card.Meta>
              <Card.Description>
                <XTaskResultDisplay id={taskDisplay.tasks[taskDisplay.active].id} />
              </Card.Description>
            </Card.Content>
          </Card> : ResultCards(jobs)
          }

          {/* {!taskDisplay ? ResultCards(jobs) : ResultDisplay()} */}
        </XBoundary>
      </XBoundary>
    </React.Fragment>
  );
};

export default XJobResults;
