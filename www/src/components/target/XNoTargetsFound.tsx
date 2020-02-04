import * as React from "react";
import { Header, Icon, Segment } from "semantic-ui-react";

const XNoTargetsFound = () => (
  <Segment placeholder>
    <Header as="h2" icon>
      <Icon name="search" />
      No Targets Found
      <Header.Subheader>
        When new targets are created, they'll be displayed here.
      </Header.Subheader>
    </Header>
  </Segment>
);

export default XNoTargetsFound;
