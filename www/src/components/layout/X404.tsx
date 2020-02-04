import * as React from "react";
import { Header, Icon, Segment } from "semantic-ui-react";

const X404 = () => (
  <Segment placeholder>
    <Header as="h2" icon>
      <Icon name="search" />
      Page Not Found
      <Header.Subheader>
        Failed to find the page you were looking for, double check the link for
        typos.
      </Header.Subheader>
    </Header>
  </Segment>
);

export default X404;
