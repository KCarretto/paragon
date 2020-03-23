import { useQuery } from "@apollo/react-hooks";
import gql from "graphql-tag";
import * as React from "react";
import { XJobCard, XNoJobsFound } from "../components/job";
import { XBoundary, XCardGroup } from "../components/layout";
import { XErrorMessage, XLoadingMessage } from "../components/messages";

export const MULTI_JOB_QUERY = gql`
  {
    jobs {
      id
      name
      staged

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
        error

        target {
          id
          name
        }

        job {
          id
          name
          staged
        }
      }
    }
  }
`;

const XMultiJobView = () => {
  const { loading, error, data: { jobs = [] } = {} } = useQuery(
    MULTI_JOB_QUERY,
    {
      pollInterval: 5000
    }
  );

  const whenLoading = (
    <XLoadingMessage title="Loading Jobs" msg="Fetching job list" />
  );
  const whenEmpty = <XNoJobsFound />;

  return (
    <React.Fragment>
      <XErrorMessage title="Error Loading Jobs" err={error} />
      <XBoundary boundary={whenLoading} show={!loading}>
        <XBoundary boundary={whenEmpty} show={jobs && jobs.length > 0}>
          <XCardGroup>
            {jobs && jobs.map(job => <XJobCard key={job.id} {...job} />)}
          </XCardGroup>
        </XBoundary>
      </XBoundary>
    </React.Fragment>
  );
};

export default XMultiJobView;
