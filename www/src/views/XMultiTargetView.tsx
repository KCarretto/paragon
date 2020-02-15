import { useQuery } from "@apollo/react-hooks";
import gql from "graphql-tag";
import * as React from "react";
import { XBoundary, XCardGroup } from "../components/layout";
import { XErrorMessage, XLoadingMessage } from "../components/messages";
import { XNoTargetsFound } from "../components/target";
import XTargetCard from "../components/target/XTargetCard";
import { Target } from "../graphql/models";

export const MULTI_TARGET_QUERY = gql`
  {
    targets {
      id
      name
      primaryIP
      lastSeen

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

        job {
          id
          name
        }
      }
    }
  }
`;

type MultiTargetQueryResponse = {
  targets: Target[];
};

const XMultiTargetView = () => {
  const { loading, error, data: { targets = [] } = {} } = useQuery<
    MultiTargetQueryResponse
  >(MULTI_TARGET_QUERY, {
    pollInterval: 5000
  });

  const whenLoading = (
    <XLoadingMessage title="Loading Targets" msg="Fetching target list" />
  );
  const whenEmpty = <XNoTargetsFound />;

  return (
    <React.Fragment>
      <XErrorMessage title="Error Loading Targets" err={error} />
      <XBoundary boundary={whenLoading} show={!loading}>
        <XBoundary boundary={whenEmpty} show={targets.length > 0}>
          <XCardGroup>
            {targets.map(target => (
              <XTargetCard key={target.id} {...target} />
            ))}
          </XCardGroup>
        </XBoundary>
      </XBoundary>
    </React.Fragment>
  );
};

export default XMultiTargetView;
