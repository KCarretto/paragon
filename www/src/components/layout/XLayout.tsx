import * as React from "react";
import { FunctionComponent } from "react";
import { BrowserRouter as Router, Switch } from "react-router-dom";
import { RouteConfig } from "../../config/routes";
import XSidebar from "./XSidebar";

type LayoutProps = {
  routeMap: RouteConfig[];
  className: string;
};

const XLayout: FunctionComponent<LayoutProps> = props => (
  <Router>
    <Switch>
      <XSidebar routeMap={props.routeMap}>{props.children}</XSidebar>
    </Switch>
  </Router>
);

export default XLayout;
