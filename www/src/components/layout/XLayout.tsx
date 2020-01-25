import * as React from "react";
import { FunctionComponent } from "react";
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";
import XLogin from '../../views/XLogin';
import XSidebar from "./XSidebar";

type LayoutProps = {
  routeMap: Route[];
  className: string;
};

const XLayout: FunctionComponent<LayoutProps> = props => (
  <Router>
    <Switch>
      <Route
        path="/login"
        component={XLogin}
      />
      <Route
        path="/"
      >
        <XSidebar routeMap={props.routeMap}>{props.children}</XSidebar>
      </Route>
    </Switch>
  </Router>
);

export default XLayout;
