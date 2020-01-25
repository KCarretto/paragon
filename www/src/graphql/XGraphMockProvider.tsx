import { MockedProvider } from "@apollo/react-testing";
import * as React from "react";
import { mocks } from ".";

const XGraphMockProvider = props => (
  <MockedProvider mocks={mocks} addTypename={false}>
    {props.children}
  </MockedProvider>
);

export default XGraphMockProvider;
