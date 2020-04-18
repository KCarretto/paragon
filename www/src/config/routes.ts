import { SemanticICONS } from "semantic-ui-react/dist/commonjs/generic";

export type RouteConfig = {
  title: string;
  link: string;
  help: string;
  icon: SemanticICONS;
};

export const Routes: RouteConfig[] = [
  {
    title: "Events",
    link: "/event_feed",
    icon: "feed",
    help: "Recent Events",
  },
  {
    title: "Run",
    link: "/run",
    icon: "code",
    help: "Create and execute jobs",
  },
  {
    title: "Targets",
    link: "/targets",
    icon: "desktop",
    help: "View targets",
  },
  {
    title: "Jobs",
    link: "/jobs",
    icon: "cubes",
    help: "View jobs",
  },
  {
    title: "Credentials",
    link: "/credentials",
    icon: "key",
    help: "View and delete credentials",
  },
  {
    title: "Tags",
    link: "/tags",
    icon: "tags",
    help: "View and manage tags",
  },
  {
    title: "Files",
    link: "/files",
    icon: "gift",
    help: "View, upload, and download files from the CDN",
  },
  {
    title: "Profile",
    link: "/profile",
    icon: "user secret",
    help: "Your profile information",
  }
];
