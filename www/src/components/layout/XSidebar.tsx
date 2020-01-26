import Cookies from "js-cookie";
import * as React from "react";
import { FunctionComponent } from "react";
import { HotKeys } from "react-hotkeys";
import { Link, Redirect } from "react-router-dom";
import { Icon, Menu, Sidebar } from "semantic-ui-react";
import { RouteConfig } from "../../config/routes";
import { XJobQueueModal } from "../job";
import "./index.css";

type SidebarProps = {
  routeMap: RouteConfig[];
};

const XSidebar: FunctionComponent<SidebarProps> = props => {
  let userId = Cookies.get("pg-userid");
  let modal = <span />;
  const createJob = React.useCallback(() => {
    modal = <XJobQueueModal openOnStart={true} />;
  }, []);
  if (!userId) {
    return <Redirect to="/login" />;
  }

  const handlers = {
    CREATE_JOB: createJob
  };
  return (
    <Sidebar.Pushable className="XLayout">
      <Sidebar
        as={Menu}
        icon="labeled"
        animation="push"
        direction="left"
        visible
        vertical
        inverted
        width="thin"
        className="XSidebar"
      >
        {props.routeMap
          ? props.routeMap.map((value: RouteConfig, index: number) => {
              return (
                <Menu.Item as={Link} to={value.link} key={index}>
                  {value.icon}
                  {value.title}
                </Menu.Item>
              );
            })
          : []}
        <Menu.Item
          href="https://github.com/kcarretto/paragon/issues/new"
          target="_blank"
        >
          <Icon name="bug" />
          Bug
        </Menu.Item>
      </Sidebar>
      <Sidebar.Pusher className="XContent">
        <HotKeys handlers={handlers}>
          {props.children}
          {modal}
        </HotKeys>
      </Sidebar.Pusher>
    </Sidebar.Pushable>
  );
};
export default XSidebar;
