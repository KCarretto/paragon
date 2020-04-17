import * as React from "react";
import { Header, Icon, Segment } from "semantic-ui-react";

const XNoCredentialsFound = () => (
  <Segment placeholder>
    <Header as="h2" icon>
      <Icon name="search" />
      No Credentials Found
      <Header.Subheader>
        When new credentials are created, they'll be displayed here.
      </Header.Subheader>
    </Header>
  </Segment>
);

export default XNoCredentialsFound;
