import moment from "moment";
import * as React from "react";
import { Container, Feed, Segment } from "semantic-ui-react";
import { default as XEvent, EventKind } from "../components/event/XEventKind";

const XEventFeedView = () => (
  <Container fluid style={{ padding: "20px" }}>
    <Segment raised>
      <Feed size="large">
        <XEvent
          event={{
            id: "123",
            kind: EventKind.JobCompleted,
            creationTime: moment().format("YYYY-MM-DDTHH:mm:ssZ"),
            owner: { id: "1234", name: "Kyle" },
            job: { id: "12345", name: "Initial Deployment" }
          }}
          kind={EventKind.JobCompleted}
        />
        <XEvent
          event={{
            id: "124",
            kind: EventKind.JobCompleted,
            creationTime: moment().format("YYYY-MM-DDTHH:mm:ssZ"),
            owner: { id: "1235", name: "Nick" },
            job: { id: "12346", name: "Initial Setup" }
          }}
          kind={EventKind.JobCompleted}
        />
      </Feed>
    </Segment>
  </Container>
);

export default XEventFeedView;
