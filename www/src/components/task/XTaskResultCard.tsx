import moment from "moment";
import * as React from "react";
import { useState } from "react";
import { Button, Card, Icon, Label } from "semantic-ui-react";
import { XTaskStatus } from ".";
import { Task } from "../../graphql/models";

const XTaskResultCard: React.FC<{
  targetName: string;
  tasks: Task[];
  onShowResult: (active: number, tasks: Task[]) => void;
}> = ({ targetName, tasks, onShowResult }) => {
  const [active, setActive] = useState<number>(0);

  return (
    <Card>
      <Card.Content>
        <Card.Header>
          {/* <Card.Header as={Link} to={getLink(task, display)}> */}
          <Icon
            floated="right"
            size="large"
            {...new XTaskStatus().getStatus(tasks[active]).icon}
          />
          {targetName}{" "}
          <Label>
            {active !== 0
              ? `Version ${tasks.length - active}`
              : `Latest (v${tasks.length})`}
          </Label>
        </Card.Header>
        <Card.Meta textAlign="center" style={{ verticalAlign: "middle" }}>
          {moment(new XTaskStatus().getTimestamp(tasks[active])).fromNow()}
          <Button
            basic
            animated
            color="blue"
            size="small"
            floated="right"
            onClick={() => onShowResult(active, tasks)}
            // as={Link}
            // to={"/tasks/" + task.id}
          >
            <Button.Content visible>View Results</Button.Content>
            <Button.Content hidden>
              <Icon name="arrow right" />
            </Button.Content>
          </Button>
        </Card.Meta>
      </Card.Content>
    </Card>
  );
};

export default XTaskResultCard;
