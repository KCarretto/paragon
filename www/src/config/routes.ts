import { SemanticICONS } from "semantic-ui-react/dist/commonjs/generic";

export type RouteConfig = {
  title: string;
  link: string;
  icon: SemanticICONS;
};

export const Routes: RouteConfig[] = [
  {
    title: "Events",
    link: "/event_feed",
    icon: "road"
  },
  {
    title: "Run",
    link: "/run",
    icon: "code"
  },
  {
    title: "Targets",
    link: "/targets",
    icon: "desktop"
  },
  {
    title: "Jobs",
    link: "/jobs",
    icon: "cubes"
  },
  {
    title: "Tags",
    link: "/tags",
    icon: "tags"
  },
  {
    title: "Files",
    link: "/files",
    icon: "gift"
  },
  {
    title: "Profile",
    link: "/profile",
    icon: "user secret"
  }
];
