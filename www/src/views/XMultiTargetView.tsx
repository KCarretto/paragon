import { useQuery } from "@apollo/react-hooks";
import gql from "graphql-tag";
import * as React from "react";
import { useState } from "react";
import { Container, Menu, Responsive } from "semantic-ui-react";
import { XTargetTypeahead } from "../components/form";
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
          staged
        }
      }
    }
  }
`;

type MultiTargetQueryResponse = {
  targets: Target[];
};

const XMultiTargetView = () => {
  const [selected, setSelected] = useState<string[]>([]);

  const { loading, error, data: { targets = [] } = {} } = useQuery<
    MultiTargetQueryResponse
  >(MULTI_TARGET_QUERY, {
    pollInterval: 5000,
  });

  const whenLoading = (
    <XLoadingMessage title="Loading Targets" msg="Fetching target list" />
  );
  const whenEmpty = <XNoTargetsFound />;

  return (
    <React.Fragment>
      <Responsive minWidth={350}>
        <Menu
          // secondary
          fixed="top"
          borderless
          fluid
          style={{
            zIndex: 50,
            margin: "0px 0px 25px 0px",
            maxHeight: "50px",
            height: "50px",
          }}
        >
          <Menu.Item
            position="right"
            style={{
              minWidth: "300px",
              width: "100%",
              paddingLeft: "200px",
            }}
          >
            <XTargetTypeahead
              labeled
              input={{ fluid: true, icon: "desktop", label: "Filter" }}
              onChange={(e, { value }) => setSelected(value || [])}
            />
          </Menu.Item>
        </Menu>
        <div style={{ marginBottom: "60px" }} />
      </Responsive>
      <XErrorMessage title="Error Loading Targets" err={error} />
      <XBoundary boundary={whenLoading} show={!loading}>
        <XBoundary boundary={whenEmpty} show={targets && targets.length > 0}>
          <Container fluid style={{ marginTop: "5px", padding: "10px" }}>
            <XCardGroup>
              {targets &&
                targets
                  .filter(
                    (target) =>
                      !selected ||
                      selected.length < 1 ||
                      selected.includes(target.id)
                  )
                  .map((target) => <XTargetCard key={target.id} {...target} />)}
            </XCardGroup>
          </Container>
        </XBoundary>
      </XBoundary>
    </React.Fragment>
  );
};

export default XMultiTargetView;
