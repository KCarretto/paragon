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
  <XSidebar
    routeMap={props.routeMap}
    userID={props.userID}
    isAdmin={props.isAdmin}
  >
    {props.children}
  </XSidebar>
);
// <Router>
//   <Switch>
//     <Route
//       path="/login"
//       render={routeProps =>
//         !authenticated || (!activated && !admin) ? (
//           <XLogin pending={authenticated} />
//         ) : (
//           <Redirect to="/" />
//         )
//       }
//     />
//     <Route
//       path="/"
//       render={routeProps =>
//         authenticated && (activated || admin) ? (
//           <XSidebar
//             routeMap={props.routeMap}
//             userID={userID}
//             isActivated={activated}
//             isAdmin={admin}
//           >
//             {props.children}
//           </XSidebar>
//         ) : (
//           <Redirect to="/login" />
//         )
//       }
//     />
//   </Switch>
// </Router>
//   );
// };

export default XLayout;
