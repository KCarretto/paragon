import { useQuery } from "@apollo/react-hooks";
import gql from "graphql-tag";
import * as React from "react";
import { FunctionComponent, useState } from "react";
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";
import { RouteConfig } from "../../config/routes";
import { User } from "../../graphql/models";
import XLogin from "../../views/XLogin";
import XPendingActivation from "../../views/XPendingActivation";
import { XLoadingMessage } from "../messages";
import XSidebar from "./XSidebar";

const WHOAMI_QUERY = gql`
  query WhoAmI {
    me {
      id
      activated
      isAdmin
    }
  }
`;

type WhoAmIResult = {
  me: User;
};

type LayoutProps = {
  routeMap: RouteConfig[];
  className: string;
};

const XLayout: FunctionComponent<LayoutProps> = props => {
  const [userID, setUserID] = useState<string>(null);
  const [activated, setActivated] = useState(false);
  const [admin, setAdmin] = useState(false);

  const { loading } = useQuery<WhoAmIResult>(WHOAMI_QUERY, {
    fetchPolicy: "no-cache",
    onCompleted: data => {
      setUserID(data.me.id);
      setActivated(data.me.activated);
      setAdmin(data.me.isAdmin);
    },
    onError: err => {
      setUserID(null);
      setActivated(false);
      setAdmin(false);
    }
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
        <Route path="/login">
          <XLogin userID={userID} isActivated={activated} isAdmin={admin} />
        </Route>
        <Route path="/login/pending">
          <XPendingActivation isActivated={activated} isAdmin={admin} />
        </Route>

        <Route path="/">
          <XSidebar
            routeMap={props.routeMap}
            userID={userID}
            isActivated={activated}
            isAdmin={admin}
          >
            {props.children}
          </XSidebar>
        </Route>
      </Switch>
    </Router>
  );
};

export default XLayout;
