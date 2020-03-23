import * as React from "react";
import { Header, Icon, Segment } from "semantic-ui-react";

type Params = {
  output?: string;
};

const XTaskOutput = ({ output }: Params) => (
  <React.Fragment>
    <Header inverted size="large" attached="top" style={{ marginTop: "10px" }}>
      <Icon name="envelope open outline" />
      <Header.Content>Output</Header.Content>
    </Header>
    <Segment raised attached style={{ overflow: 'auto' }}>
      <pre>{output !== null ? output : "No Output Available"}</pre>
    </Segment>
  </React.Fragment>
);

export default XTaskOutput;
