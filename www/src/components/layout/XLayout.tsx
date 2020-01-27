import { useQuery } from "@apollo/react-hooks";
import gql from "graphql-tag";
import * as React from "react";
import { FunctionComponent } from "react";
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";
import { RouteConfig } from "../../config/routes";
import { User } from "../../graphql/models";
import XLogin from "../../views/XLogin";
import { XLoadingMessage } from "../messages";
import XSidebar from "./XSidebar";

const WHOAMI_QUERY = gql`
  query WhoAmI {
    me {
      id
    }
  }
`;

type WhoAmIResult = {
  data: User;
};

type LayoutProps = {
  routeMap: RouteConfig[];
  className: string;
};

const XLayout: FunctionComponent<LayoutProps> = props => {
  const { loading } = useQuery<WhoAmIResult>(WHOAMI_QUERY, {
    fetchPolicy: "no-cache"
  });

  if (loading) {
    return (
      <XLoadingMessage
        title="Loading Paragon"
        msg="Querying identity..."
        hidden={false}
      />
    );
  }

  return (
    <Router>
      <Switch>
        <Route path="/login" component={XLogin} />
        <Route path="/">
          <XSidebar routeMap={props.routeMap}>{props.children}</XSidebar>
        </Route>
      </Switch>
    </Router>
  );
};

export default XLayout;
