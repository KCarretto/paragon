import moment from "moment";
import * as React from "react";
import { FunctionComponent } from "react";
import { Link } from "react-router-dom";
import { Button, Card, Icon } from "semantic-ui-react";
import { XTaskStatus } from ".";
import { Task } from "../../graphql/models";

export type XTaskCardProps = {
  task: Task;
  display: XTaskCardDisplayType;
};

export enum XTaskCardDisplayType {
  JOB = 1,
  TARGET = 2
}

const getLink = (task: Task, displayType: XTaskCardDisplayType): string => {
  switch (displayType) {
    case XTaskCardDisplayType.TARGET:
      return "/targets/" + (task.target ? task.target.id : "0");
    case XTaskCardDisplayType.JOB:
      return "/jobs/" + (task.job ? task.job.id : "0");
  }
  console.error("Failed to get link for XTaskCard");
};

const getHeader = (task: Task, displayType: XTaskCardDisplayType): string => {
  switch (displayType) {
    case XTaskCardDisplayType.TARGET:
      return task.target ? task.target.name : "Untitled Target";
    case XTaskCardDisplayType.JOB:
      return task.job ? task.job.name : "Untitled Job";
  }
};

const XTaskCard: FunctionComponent<XTaskCardProps> = ({ task, display }) => {
  const target = task.target || { id: 0, name: "Untitled Target" };
  return (
    <Card>
      <Card.Content>
        <Card.Header as={Link} to={getLink(task, display)}>
          <Icon
            floated="right"
            size="large"
            {...new XTaskStatus().getStatus(task).icon}
          />
          {getHeader(task, display)}
        </Card.Header>
        <Card.Meta textAlign="center" style={{ verticalAlign: "middle" }}>
          {moment(new XTaskStatus().getTimestamp(task)).fromNow()}
          <Button
            basic
            animated
            color="blue"
            size="small"
            floated="right"
            as={Link}
            to={"/tasks/" + task.id}
          >
            <Button.Content visible>View Task</Button.Content>
            <Button.Content hidden>
              <Icon name="arrow right" />
            </Button.Content>
          </Button>
        </Card.Meta>
      </Card.Content>
    </Card>
  );
};

export default XTaskCard;
