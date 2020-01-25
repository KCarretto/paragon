import { Icon } from "semantic-ui-react";

export type Route = {
  title: String;
  link: String;
  icon: Icon;
};

const routes: Route[] = [
  {
    title: "News Feed",
    link: "/news_feed",
    icon: new Icon({name: "newspaper"})
  },
  {
    title: "Targets",
    link: "/targets",
    icon: new Icon({name: "desktop"})
  },
  {
    title: "Jobs",
    link: "/jobs",
    icon: new Icon({name: "cubes"})
  },
  {
    title: "Tags",
    link: "/tags",
    icon: new Icon({name: "tags"})
  },
  {
    title: "Files",
    link: "/files",
    icon: new Icon({name: "gift"})
  },
  {
    title: "Profile",
    link: "/profile",
    icon: new Icon({name: "user secret"})
  }
];

export default routes;
