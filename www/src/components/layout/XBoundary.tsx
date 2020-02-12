import * as React from "react";
import { FunctionComponent } from "react";

type BoundaryProps = {
  show?: boolean;
  boundary: React.ReactNode;
};

const XBoundary: FunctionComponent<BoundaryProps> = props => (
  <React.Fragment>
    {props.show ? props.children : props.boundary}
  </React.Fragment>
);

export default XBoundary;
