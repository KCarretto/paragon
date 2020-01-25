import { ApolloError } from "apollo-client/errors/ApolloError";
import * as React from "react";
import { Message } from "semantic-ui-react";

export type ErrorMessageParams = {
  title: string;
  err: ApolloError;
};

export default ({ title, err }: ErrorMessageParams) => (
  <Message negative hidden={!err ? true : false}>
    <Message.Header>{title}</Message.Header>
    {err && err.message ? err.message : ""}
  </Message>
);
