import * as React from "react";
import { Header, Icon, Segment } from "semantic-ui-react";

const XNoJobsFound = () => (
  <Segment placeholder>
    <Header as="h2" icon>
      <Icon name="search" />
      No Jobs Found
      <Header.Subheader>
        When new jobs are queued, they'll be displayed here.
      </Header.Subheader>
    </Header>
  </Segment>
);

export default XNoJobsFound;
