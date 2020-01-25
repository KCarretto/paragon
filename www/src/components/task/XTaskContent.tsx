import * as React from "react";
import { Header, Icon, Segment } from "semantic-ui-react";

export default (content: String) => (
  <div>
    <Header inverted size="large" attached="top" style={{ marginTop: "10px" }}>
      <Icon name="code" />
      <Header.Content>Content</Header.Content>
    </Header>
    <Segment raised attached>
      <pre>{content || "No Content Available"}</pre>
    </Segment>
  </div>
);
