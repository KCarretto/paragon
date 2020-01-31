import { useQuery } from "@apollo/react-hooks";
import gql from "graphql-tag";
import * as React from "react";
import { Container, Feed, Loader, Segment } from "semantic-ui-react";
import { default as XEvent, EventKind } from "../components/event/XEventKind";
import XNoEventsFound from "../components/event/XNoEventsFound";
import { XErrorMessage } from "../components/messages";
import { Event } from "../graphql/models";

export const EVENT_FEED_QUERY = gql`
  {
    events {
      id
      creationTime
      kind

      job {
        id
        name
      }
      file {
        id
        name
        size
      }
      credential {
        id
        secret
        principal
        kind
      }
      link {
        id
        alias
        expirationTime
        clicks
      }
      tag {
        id
        name
      }
      target {
        id
        name
      }
      task {
        id
        claimTime
        execStartTime
        execStopTime
      }
      user {
        id
        name
        isActivated
        isAdmin
      }
      service {
        id
        name
        isActivated
      }
      event {
        id
        kind
      }
      likers {
        id
      }
      owner {
        id
        name
        isActivated
        isAdmin
      }
      svcOwner {
        id
        name
        isActivated
      }
    }
  }
`;

export type EvnetFeedResponse = {
  events: Event[];
};

const XEventFeedView = () => {
  const { called, loading, error, data } = useQuery<EvnetFeedResponse>(
    EVENT_FEED_QUERY
  );

  if (!data || !data.events || data.events.length < 1) {
    return <XNoEventsFound />;
  }

  return (
    <Container fluid style={{ padding: "20px" }}>
      <Loader disabled={!called || !loading} />
      <Segment raised>
        <Feed size="large">
          {data.events.map(e => {
            console.log(e.kind);
            console.log(EventKind[e.kind]);
            return <XEvent event={e} kind={EventKind[e.kind]} />;
          })}
        </Feed>
      </Segment>
      <XErrorMessage title="Error Loading Feed" err={error} />
    </Container>
  );
};

export default XEventFeedView;
