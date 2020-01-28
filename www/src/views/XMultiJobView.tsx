import { useQuery } from "@apollo/react-hooks";
import gql from "graphql-tag";
import * as React from "react";
import { Card, Container, Loader, Menu } from "semantic-ui-react";
import { XJobCard, XJobQueueModal, XNoJobsFound } from "../components/job";
import { XErrorMessage } from "../components/messages";

export const MULTI_JOB_QUERY = gql`
  {
    jobs {
      id
      name

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

        job {
          id
          name
        }
      }
    }
  }
`;

const XMultiJobView = () => {
  const { called, loading, error, data } = useQuery(MULTI_JOB_QUERY, {
    pollInterval: 5000
  });

  const showCards = () => {
    if (!data || !data.jobs || data.jobs.length < 1) {
      return <XNoJobsFound />;
    }
    return (
      <Card.Group centered itemsPerRow={4}>
        {data.jobs.map(job => (
          <XJobCard key={job.id} {...job} />
        ))}
      </Card.Group>
    );
  };

  return (
    <Container style={{ padding: "10px" }}>
      <Menu secondary>
        <Menu.Item position="right">
          <XJobQueueModal />
        </Menu.Item>
      </Menu>
      <Container fluid style={{ padding: "20px" }}>
        <Loader disabled={!called || !loading} />
        <XErrorMessage title="Error Loading Jobs" err={error} />
        {showCards()}
      </Container>
    </Container>
  );
};

export default XMultiJobView;
