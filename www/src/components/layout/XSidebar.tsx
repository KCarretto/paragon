import * as React from "react";
import { FunctionComponent } from "react";
import { HotKeys } from "react-hotkeys";
import { Link } from "react-router-dom";
import { Container, Icon, Menu, Responsive, Sidebar } from "semantic-ui-react";
import { RouteConfig } from "../../config/routes";
import { XJobQueueModal } from "../job";
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
            <Menu.Item as={Link} to={value.link} key={index}>
              {value.icon}
              <Responsive minWidth={800}>{value.title}</Responsive>
            </Menu.Item>
          );
        })
      : []}
    {props.isAdmin ? (
      <Menu.Item as={Link} to="/admin">
        <Icon name="chess rook" />
        <Responsive minWidth={800}>Admin</Responsive>
      </Menu.Item>
    ) : (
      <span />
    )}
    <Menu.Item
      href="https://github.com/kcarretto/paragon/issues/new"
      target="_blank"
    >
      <Icon name="bug" />
      <Responsive minWidth={800}>Bug</Responsive>
    </Menu.Item>
  </React.Fragment>
);

const mobileSidebar = (props: SidebarProps) => (
  <Responsive maxWidth={799}>
    <Sidebar
      animation="push"
      direction="left"
      visible
      inverted
      as={Menu}
      icon
      vertical
      compact
      className="XSidebar"
    >
      {/* <Menu icon vertical compact inverted className="XSidebar"> */}
      {getMenu(props)}
      {/* </Menu> */}
    </Sidebar>
  </Responsive>
);

const desktopSidebar = (props: SidebarProps) => (
  <Responsive minWidth={800}>
    <Sidebar
      as={Menu}
      icon="labeled"
      animation="push"
      direction="left"
      visible
      vertical
      compact
      inverted
      // width="very thin"
      className="XSidebar"
    >
      {getMenu(props)}
    </Sidebar>
  </Responsive>
);

const XSidebar: FunctionComponent<SidebarProps> = props => {
  // let userId = Cookies.get("pg-userid");
  let modal = <span />;
  const createJob = React.useCallback(() => {
    modal = <XJobQueueModal openOnStart={true} />;
  }, []);

  const handlers = {
    CREATE_JOB: createJob
  };
  return (
    <Sidebar.Pushable className="XLayout">
      {mobileSidebar(props)}
      {desktopSidebar(props)}
      <Sidebar.Pusher>
        <HotKeys handlers={handlers}>
          <div className="XContent">
            <Responsive maxWidth={799}>
              <Container style={{ paddingLeft: "5vw" }}>
                {props.children}
                {modal}
              </Container>
            </Responsive>
            <Responsive minWidth={800}>
              {props.children}
              {modal}
            </Responsive>
          </div>

          {/* </Container> */}
        </HotKeys>
      </Sidebar.Pusher>
    </Sidebar.Pushable>
  );
};
export default XSidebar;
