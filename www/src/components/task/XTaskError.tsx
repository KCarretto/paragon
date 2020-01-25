import * as React from "react";
import { Header, Icon, Segment } from "semantic-ui-react";

type Params = {
  error: string;
};

const XTaskError = ({ error }: Params) => {
  if (!error) {
    return <span />;
  }

  return (
    <React.Fragment>
      <Header
        inverted
        size="large"
        attached="top"
        color="red"
        style={{ marginTop: "10px" }}
      >
        <Icon name="warning circle" />
        <Header.Content>Error</Header.Content>
      </Header>
      <Segment raised attached>
        <pre>{error}</pre>
      </Segment>
    </React.Fragment>
  );
};

export default XTaskError;
