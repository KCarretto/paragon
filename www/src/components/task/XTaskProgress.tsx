import * as React from "react";
import { Icon, Step } from "semantic-ui-react";

const XTaskProgress = () => (
  <Step.Group size="mini" vertical>
    <Step completed>
      <Icon name="wait" size="mini" />
      <Step.Content>
        <Step.Title>Queued</Step.Title>
        <Step.Description>5m ago</Step.Description>
      </Step.Content>
    </Step>
    <Step completed>
      <Icon name="send" size="mini" />
      <Step.Content>
        <Step.Title>Claimed</Step.Title>
        <Step.Description>1m ago</Step.Description>
      </Step.Content>
    </Step>
    <Step active>
      <Icon name="bolt" size="mini" />
      <Step.Content>
        <Step.Title>Executing</Step.Title>
        <Step.Description>Started 1m ago</Step.Description>
      </Step.Content>
    </Step>
    <Step disabled>
      <Icon name="check circle" size="mini" />
      <Step.Content>
        <Step.Title>Completed</Step.Title>
      </Step.Content>
    </Step>
  </Step.Group>
);
export default XTaskProgress;
