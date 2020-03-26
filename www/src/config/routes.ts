import * as React from "react";
import { Icon, IconProps } from "semantic-ui-react";

export type RouteConfig = {
  title: string;
  link: string;
  icon: React.CElement<IconProps, Icon>;
};

export const Routes: RouteConfig[] = [
  {
    title: "Event Feed",
    link: "/event_feed",
    icon: React.createElement(Icon, { name: "road" })
  },
  {
    title: "Run",
    link: "/run",
    icon: React.createElement(Icon, { name: "code" })
  },
  {
    title: "Targets",
    link: "/targets",
    icon: React.createElement(Icon, { name: "desktop" })
  },
  {
    title: "Jobs",
    link: "/jobs",
    icon: React.createElement(Icon, { name: "cubes" })
  },
  {
    title: "Tags",
    link: "/tags",
    icon: React.createElement(Icon, { name: "tags" })
  },
  {
    title: "Files",
    link: "/files",
    icon: React.createElement(Icon, { name: "gift" })
  },
  {
    title: "Profile",
    link: "/profile",
    icon: React.createElement(Icon, { name: "user secret" })
  }
];
