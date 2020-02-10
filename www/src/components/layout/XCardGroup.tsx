import * as React from "react";
import { FunctionComponent } from "react";
import { Card } from "semantic-ui-react";

const XCardGroup: FunctionComponent = props => (
  <Card.Group centered stackable itemsPerRow={4}>
    {props.children}
  </Card.Group>
);

export default XCardGroup;
