import * as React from "react";
import { FunctionComponent } from "react";
import { Link } from "react-router-dom";
import { Icon, Menu, Sidebar } from "semantic-ui-react";
import { Route } from "../../config/routes";
import "./index.css";

type SidebarProps = {
  routeMap: Route[];
};

const XSidebar: FunctionComponent<SidebarProps> = props => (
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
        ? props.routeMap.map((value: Route, index: number) => {
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
    <Sidebar.Pusher className="XContent">{props.children}</Sidebar.Pusher>
  </Sidebar.Pushable>
);

export default XSidebar;
