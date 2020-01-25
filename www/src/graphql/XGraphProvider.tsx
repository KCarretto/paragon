import { InMemoryCache } from "apollo-cache-inmemory";
import { ApolloClient } from "apollo-client";
import { createHttpLink } from "apollo-link-http";
import * as React from "react";
import { ApolloProvider } from "react-apollo";

const httpLink = createHttpLink({
  uri: window.location.origin + "/graphql"
});

const client = new ApolloClient({
  link: httpLink,
  cache: new InMemoryCache()
});

const XGraphProvider = props => (
  <ApolloProvider client={client}>{props.children}</ApolloProvider>
);

export default XGraphProvider;
