import { useQuery } from "@apollo/react-hooks";
import gql from "graphql-tag";
import * as React from "react";
import { Loader } from "semantic-ui-react";
import { XJobCard, XNoJobsFound } from "../components/job";
import { XCardGroup } from "../components/layout";
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

        target {
          id
          name
        }

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
      <XCardGroup>
        {data.jobs.map(job => (
          <XJobCard key={job.id} {...job} />
        ))}
      </XCardGroup>
    );
  };

  return (
    <React.Fragment>
      <Loader disabled={!called || !loading} />
      <XErrorMessage title="Error Loading Jobs" err={error} />
      {showCards()}
    </React.Fragment>
  );
};

export default XMultiJobView;
