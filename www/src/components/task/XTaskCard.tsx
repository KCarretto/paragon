import moment from "moment";
import * as React from "react";
import { Link } from "react-router-dom";
import { Button, Card, Icon } from "semantic-ui-react";
import { XTaskStatus } from ".";
import { Task } from "../../graphql/models";

const XTaskCard = (task: Task) => {
  const target = task.target || { id: 0, name: "Untitled Target" };
  return (
    <Card>
      <Card.Content>
        <Card.Header as={Link} to={"/targets/" + target.id}>
          <Icon
            floated="right"
            size="large"
            {...new XTaskStatus().getStatus(task).icon}
          />
          {target.name}
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
