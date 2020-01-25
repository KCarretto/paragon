import * as React from "react";
import { Header, Icon, Segment } from "semantic-ui-react";

export default (output?: String) => (
  <div>
    <Header inverted size="large" attached="top" style={{ marginTop: "10px" }}>
      <Icon name="envelope open outline" />
      <Header.Content>Output</Header.Content>
    </Header>
    <Segment raised attached>
      <pre>{output || "No Output Available"}</pre>
    </Segment>
  </div>
);
