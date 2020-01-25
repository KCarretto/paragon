import moment from "moment";
import * as React from "react";
import { Card, Label } from "semantic-ui-react";
import { Target } from "../../graphql/models";
import { XTags } from "../tag";
import { XTaskSummary } from "../task";

const XTargetCard = (t: Target) => (
  <Card fluid>
    <Card.Content>
      <Card.Header href={"/targets/" + t.id}>{t.name} </Card.Header>
      {!t.lastSeen ||
      moment(t.lastSeen).isBefore(moment().subtract(5, "minutes")) ? (
        <Label corner="right" size="large" icon="times circle" color="red" />
      ) : (
        <Label corner="right" size="large" icon="check circle" color="green" />
      )}
      <Card.Meta>
        {t.lastSeen ? moment(t.lastSeen).fromNow() : "Never"}
      </Card.Meta>
      <XTaskSummary tasks={t.tasks} />
    </Card.Content>
    <Card.Content extra>
      <XTags tags={t.tags} defaultText="None" />
    </Card.Content>
  </Card>
);

export default XTargetCard;
