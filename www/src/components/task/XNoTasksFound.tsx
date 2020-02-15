import * as React from "react";
import { Header, Icon, Segment } from "semantic-ui-react";

const XNoTasksFound = () => (
  <Segment placeholder>
    <Header as="h2" icon>
      <Icon name="search" />
      No Tasks Found
      <Header.Subheader>
        When new tasks are created, they'll be displayed here.
      </Header.Subheader>
    </Header>
  </Segment>
);

export default XNoTasksFound;
