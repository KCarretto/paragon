import * as React from "react";
import { Header, Icon, Segment } from "semantic-ui-react";

const XNoTagsFound = () => (
  <Segment placeholder>
    <Header as="h2" icon>
      <Icon name="search" />
      No Tags Found
      <Header.Subheader>
        When new tags are created, they'll be displayed here.
      </Header.Subheader>
    </Header>
  </Segment>
);

export default XNoTagsFound;
