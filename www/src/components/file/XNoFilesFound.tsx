import * as React from "react";
import { Header, Icon, Segment } from "semantic-ui-react";

const XNoFilesFound = () => (
  <Segment placeholder>
    <Header as="h2" icon>
      <Icon name="search" />
      No Files Found
      <Header.Subheader>
        When new files are uploaded, they'll be displayed here.
      </Header.Subheader>
    </Header>
  </Segment>
);

export default XNoFilesFound;
