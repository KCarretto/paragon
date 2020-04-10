import * as React from "react";
import { Button, Checkbox, Grid, Modal } from "semantic-ui-react";
import { XServiceTypeahead } from "../form";

const XJobSettingsModal: React.FC<{
  serviceTag: string;
  setServiceTag: (string) => void;
  stage: boolean;
  setStage: (boolean) => void;
}> = ({ serviceTag, setServiceTag, stage, setStage }) => (
  <Modal
    closeIcon
    trigger={<Button icon="cogs" color="teal" />}
    centered={false}
  >
    <Modal.Header>Job Settings</Modal.Header>
    <Modal.Content>
      <Grid verticalAlign="middle" stackable container columns={"equal"}>
        <Grid.Row>
          <Grid.Column>
            <XServiceTypeahead
              labeled
              value={serviceTag}
              defaultSVC={"pg-worker"}
              onChange={(e, { value }) => setServiceTag(value)}
            />
            {/* <XTagTypeahead
              labeled
              onChange={(e, { value }) => setTags(value)}
            /> */}
          </Grid.Column>
          <Grid.Column>
            <Checkbox
              label="Stage this job"
              onChange={() => setStage(!stage)}
              checked={stage}
            />
          </Grid.Column>
        </Grid.Row>
      </Grid>
    </Modal.Content>
  </Modal>
);

export default XJobSettingsModal;
