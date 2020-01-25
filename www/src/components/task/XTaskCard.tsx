import moment from "moment";
import * as React from "react";
import { Button, Card, Icon } from "semantic-ui-react";
import { XTaskStatus } from ".";
import { Task } from "../../graphql/models";

const XTaskCard = (task: Task) => {
  const target = task.target || { id: 0, name: "Untitled Target" };

  return (
    <Card>
      <Card.Content>
        <Card.Header href={"/targets/" + target.id}>
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
            href={"/tasks/" + task.id}
            color="blue"
            size="small"
            floated="right"
          >
            <Button.Content visible>View Task</Button.Content>
            <Button.Content hidden>
              <Icon name="arrow right" />
            </Button.Content>
          </Button>
        </Card.Meta>
      </Card.Content>
      {/* <Card.Content extra textAlign='center'>

            </Card.Content> */}
    </Card>
  );
};

export default XTaskCard;
