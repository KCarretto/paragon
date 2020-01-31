import * as React from "react";
import { FunctionComponent } from "react";
import { Button, Container, Menu } from "semantic-ui-react";
import { XSidebar } from ".";
import { RouteConfig } from "../../config/routes";
import { XBulkAddCredentialsModal } from "../credential";
import { XFileUploadModal } from "../file";
import { XJobQueueModal } from "../job";
import { XTagCreateModal } from "../tag";
import { XTargetCreateModal } from "../target";

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
    <Menu secondary compact fixed="top">
      <Menu.Item position="right">
        <Button.Group icon color="green">
          <XJobQueueModal />
          <XFileUploadModal button={{ color: "green", icon: "cloud upload" }} />
          <XTagCreateModal />
          <XBulkAddCredentialsModal />
          <XTargetCreateModal />
        </Button.Group>
      </Menu.Item>
    </Menu>

    <Container
      fluid
      style={{ padding: "10px", paddingTop: "50px", height: "100vh" }}
    >
      {props.children}
    </Container>
  </XSidebar>
);

export default XLayout;
