import { useQuery } from "@apollo/react-hooks";
import gql from "graphql-tag";
import * as React from "react";
import { useParams } from "react-router-dom";
import { Header, Icon } from "semantic-ui-react";
import { XJobHeader } from "../components/job";
import { XBoundary, XCardGroup } from "../components/layout";
import { XErrorMessage, XLoadingMessage } from "../components/messages";
import { XNoTasksFound, XTaskCard, XTaskContent } from "../components/task";
import { XTaskCardDisplayType } from "../components/task/XTaskCard";
import { Job } from "../graphql/models";

export const JOB_QUERY = gql`
  query Job($id: ID!) {
    job(id: $id) {
      id
      name
      content
      tags {
        id
        name
      }
      tasks {
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
  }
`;

type JobQuery = {
  job: Job;
};

const XJobView = () => {
  let { id } = useParams();

  const {
    loading,
    error,
    data: { job: { name = "", content = "", tags = [], tasks = [] } = {} } = {}
  } = useQuery<JobQuery>(JOB_QUERY, {
    variables: { id },
    pollInterval: 5000
  });

  const whenLoading = (
    <XLoadingMessage title="Loading Job" msg="Fetching job information..." />
  );
  const whenEmpty = <XNoTasksFound />;

  return (
    <React.Fragment>
      <XJobHeader name={name} tags={tags} />

      <XErrorMessage title="Error Loading Job" err={error} />

      <XBoundary boundary={whenLoading} show={!loading}>
        <XTaskContent content={content} />
        <Header size="large" block inverted>
          <Icon name="tasks" />
          <Header.Content>Tasks</Header.Content>
        </Header>
        <XBoundary boundary={whenEmpty} show={tasks && tasks.length > 0}>
          <XCardGroup>
            {tasks &&
              tasks.map(task => (
                <XTaskCard
                  key={task.id}
                  display={XTaskCardDisplayType.TARGET}
                  task={task}
                />
              ))}
          </XCardGroup>
        </XBoundary>
      </XBoundary>
    </React.Fragment>
  );
};

export default XJobView;
