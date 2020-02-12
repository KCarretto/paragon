import * as React from "react";
import { Header, Icon, Segment } from "semantic-ui-react";

type Params = {
  content: string;
};

const XTaskContent = ({ content }: Params) => (
  <React.Fragment>
    <Header inverted size="large" attached="top" style={{ marginTop: "10px" }}>
      <Icon name="code" />
      <Header.Content>Content</Header.Content>
    </Header>
    <Segment raised attached style={{ overflow: "auto" }}>
      <pre>{content !== null ? content : "No Content Available"}</pre>
    </Segment>
  </React.Fragment>
);

export default XTaskContent;
