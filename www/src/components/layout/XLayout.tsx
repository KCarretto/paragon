import * as React from "react";
import { FunctionComponent } from "react";
import { XSidebar } from ".";
import { RouteConfig } from "../../config/routes";

type LayoutProps = {
  routeMap: RouteConfig[];
  userID?: string;
  isAdmin: boolean;
  className: string;
};

const XLayout: FunctionComponent<LayoutProps> = props => (
  <React.Fragment>
    <XSidebar
      routeMap={props.routeMap}
      userID={props.userID}
      isAdmin={props.isAdmin}
    />
    <div className="XContent">{props.children}</div>
  </React.Fragment>
);

export default XLayout;
