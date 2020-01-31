import { useQuery } from "@apollo/react-hooks";
import gql from "graphql-tag";
import * as React from "react";
import { Card, Container, Loader } from "semantic-ui-react";
import { XErrorMessage } from "../components/messages";
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
  const { called, loading, error, data } = useQuery<MultiTargetQueryResponse>(
    MULTI_TARGET_QUERY,
    {
      pollInterval: 5000
    }
  );

  const showCards = () => {
    if (!data || !data.targets || data.targets.length < 1) {
      return <XNoTargetsFound />;
    }
    return (
      <Card.Group centered itemsPerRow={4}>
        {data.targets.map(target => (
          <XTargetCard key={target.id} {...target} />
        ))}
      </Card.Group>
    );
  };

  return (
    <Container fluid style={{ padding: "20px" }}>
      <Loader disabled={!called || !loading} />
      <XErrorMessage title="Error Loading Targets" err={error} />
      {showCards()}
    </Container>
  );
};

// XTargetCardGroup.propTypes = {
//     targets: PropTypes.arrayOf(PropTypes.shape({
//         id: PropTypes.number.isRequired,
//         name: PropTypes.string.isRequired,
//         tags: PropTypes.arrayOf(PropTypes.string),
//     })).isRequired,
// };

export default XMultiTargetView;
