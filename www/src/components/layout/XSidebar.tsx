import * as React from "react";
import { FunctionComponent } from "react";
import { Link } from "react-router-dom";
import { Icon, Menu, Responsive } from "semantic-ui-react";
import { RouteConfig } from "../../config/routes";
import "./index.css";

type SidebarProps = {
  routeMap: RouteConfig[];
  userID?: string;
  isAdmin: boolean;
};

const getMenu = (props: SidebarProps) => (
  <React.Fragment>
    {props.routeMap
      ? props.routeMap.map((value: RouteConfig, index: number) => {
        return (
          <Menu.Item key={index} as={Link} to={value.link} link style={{ display: "flex", alignItems: "center" }}>
            <Icon fitted name={value.icon} size="big" />
            <Responsive minWidth={1800}><span style={{ marginLeft: "15px" }}><b>{value.title}</b></span></Responsive>
          </Menu.Item>
        );
      })
      : []}
    {props.isAdmin ? (
      <Menu.Item as={Link} to="/admin" link style={{ display: "flex", alignItems: "center" }}>
        <Icon fitted name="chess rook" size="big" />
        <Responsive minWidth={1800}><span style={{ marginLeft: "15px" }}><b>Admin</b></span></Responsive>
      </Menu.Item>
    ) : (
        <span />
      )}
    <Menu.Item
      link style={{ display: "flex", alignItems: "center" }}
      href="https://github.com/kcarretto/paragon/issues/new"
      target="_blank"
    >
      <Icon fitted name="bug" size="big" />
      <Responsive minWidth={1800}><span style={{ marginLeft: "15px" }}><b>Bug</b></span></Responsive>
    </Menu.Item>
  </React.Fragment>
);

const mobileSidebar = (props: SidebarProps) => (
  <Responsive maxWidth={1499}>
    {/* <Menu size="large" icon vertical compact inverted className="XSidebar"> */}
    <Menu size="large" icon vertical inverted className="XSidebar">
      {getMenu(props)}
    </Menu>
  </Responsive>
);

const desktopSidebar = (props: SidebarProps) => (
  <Responsive minWidth={1500}>
    <Menu
      size="large"
      // icon="labeled"
      vertical
      inverted
      className="XSidebar"
    >
      {getMenu(props)}
    </Menu>
  </Responsive>
);

const XSidebar: FunctionComponent<SidebarProps> = props => (
  // <Sidebar.Pushable className="XLayout">
  <React.Fragment>
    <Menu icon vertical inverted className="XSidebar">
      {getMenu(props)}
    </Menu>
  </React.Fragment>
);
//   {/* <Sidebar.Pusher>
//     <div className="XContent">
//       <Responsive maxWidth={799}>
//         <Container style={{ paddingLeft: "5vw" }}>
//           {props.children}
//         </Container>
//       </Responsive>
//       <Responsive minWidth={800}>{props.children}</Responsive>
//     </div>

//     {/* </Container> */}
//   // </Sidebar.Pusher> */}
// // </Sidebar.Pushable>
export default XSidebar;
