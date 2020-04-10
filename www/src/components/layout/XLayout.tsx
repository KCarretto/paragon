import * as React from "react";
import { FunctionComponent } from "react";
import { SemanticToastContainer } from "react-semantic-toasts";
import { XSidebar } from ".";
import { RouteConfig } from "../../config/routes";

type LayoutProps = {
  routeMap: RouteConfig[];
  userID?: string;
  isAdmin: boolean;
  className: string;
};

const XLayout: FunctionComponent<LayoutProps> = props => {
  return (
    <React.Fragment>
      <XSidebar
        routeMap={props.routeMap}
        userID={props.userID}
        isAdmin={props.isAdmin}
      />
      <div className="XContent">{props.children}</div>

      <SemanticToastContainer
        // position="bottom-right"
        animation="fade up"
        className="XToastContainer"
      />

    </React.Fragment>
  );
}

export default XLayout;
