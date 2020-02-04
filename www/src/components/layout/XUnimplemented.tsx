import React from "react";
import { Container, Header, Icon, Segment } from "semantic-ui-react";

const XUnimplemented = props => (
  <Container textAlign="center" style={{ padding: "50px" }}>
    <Segment placeholder>
      <Header as="h2" icon>
        <Icon name="wrench" />
        Under Construction
        <Header.Subheader>
          <span>
            This page has yet to be implemented, please try again later.{" "}
            <a
              href="https://github.com/KCarretto/paragon/pulls"
              target="_blank"
              rel="noopener noreferrer"
            >
              Want to help?
            </a>
          </span>
        </Header.Subheader>
      </Header>
    </Segment>
  </Container>
);

export default XUnimplemented;
