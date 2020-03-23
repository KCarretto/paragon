import { useQuery } from "@apollo/react-hooks";
import gql from "graphql-tag";
import * as React from "react";
import { Feed } from "semantic-ui-react";
import { EventKind, XEvent, XNoEventsFound } from "../components/event";
import { XBoundary } from "../components/layout";
import { XErrorMessage, XLoadingMessage } from "../components/messages";
import { Event } from "../graphql/models";

export const EVENT_FEED_QUERY = gql`
  {
    events(input: { limit: 50 }) {
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
        photoURL
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
        photoURL
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

export type EventFeedResponse = {
  events: Event[];
};

const XEventFeedView = () => {
  const { loading, error, data: { events = [] } = {} } = useQuery<
    EventFeedResponse
  >(EVENT_FEED_QUERY, {
    pollInterval: 3000
  });

  const whenLoading = (
    <XLoadingMessage title="Loading Events" msg="Loading Event Feed" />
  );
  const whenEmpty = <XNoEventsFound />;

  return (
    <React.Fragment>
      <XErrorMessage title="Error Loading Feed" err={error} />

      <XBoundary boundary={whenLoading} show={!loading}>
        <XBoundary boundary={whenEmpty} show={events && events.length > 0}>
          <Feed size="large">
            {events &&
              events.map((e, index) => {
                return (
                  <XEvent key={index} event={e} kind={EventKind[e.kind]} />
                );
              })}
          </Feed>
        </XBoundary>
      </XBoundary>
    </React.Fragment>
  );
};

export default XEventFeedView;
