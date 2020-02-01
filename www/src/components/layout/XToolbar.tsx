import * as React from "react";
import { Button, Menu } from "semantic-ui-react";
import { XBulkAddCredentialsModal } from "../credential";
import { XFileUploadModal } from "../file";
import { XJobQueueModal } from "../job";
import { XTagCreateModal } from "../tag";
import { XTargetCreateModal } from "../target";

const XToolbar = () => (
  <Menu secondary compact fixed="bottom" className="XToolbar">
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
);

export default XToolbar;
