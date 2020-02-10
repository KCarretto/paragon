import * as React from "react";
import { Header, Icon, Segment } from "semantic-ui-react";

const XNoEventsFound = () => (
  <Segment placeholder>
    <Header as="h2" icon>
      <Icon name="search" />
      No Events Found
      <Header.Subheader>
        When new events occurr, they'll be added to your feed.
      </Header.Subheader>
    </Header>
  </Segment>
);

export default XNoEventsFound;
