import * as React from "react";
import { FunctionComponent } from "react";
import { HotKeys } from "react-hotkeys";
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";
import { RouteConfig } from "../../config/routes";
import XLogin from "../../views/XLogin";
import XSidebar from "./XSidebar";

const keyMap = {
  CREATE_JOB: "command+j"
};

type LayoutProps = {
  routeMap: RouteConfig[];
  className: string;
};

const XLayout: FunctionComponent<LayoutProps> = props => (
  <Router>
    <Switch>
      <Route path="/login" component={XLogin} />
      <Route path="/">
        <HotKeys keyMap={keyMap}>
          <XSidebar routeMap={props.routeMap}>{props.children}</XSidebar>
        </HotKeys>
      </Route>
    </Switch>
  </Router>
);

export default XLayout;
