import { useMutation } from "@apollo/react-hooks";
import gql from "graphql-tag";
import * as React from "react";
import { Link } from "react-router-dom";
import { Button, Card, Icon } from "semantic-ui-react";
import { Job } from "../../graphql/models";
import { MULTI_JOB_QUERY, MULTI_TARGET_QUERY } from "../../views";
import { XTags } from "../tag";
import { XTaskSummary, XTaskSummaryDisplayType } from "../task";

export const QUEUE_JOB_MUTATION = gql`
  mutation QueueJob($id: ID!) {
    queueJob(input: { id: $id }) {
      id
    }
  }
`;

const XJobCard = (j: Job) => {
  const [queueJob, { called, loading }] = useMutation(QUEUE_JOB_MUTATION, {
    refetchQueries: [{ query: MULTI_JOB_QUERY }, { query: MULTI_TARGET_QUERY }]
  });

  const handleSubmit = () => {
    let vars = {
      id: j.id
    };

    queueJob({
      variables: vars
    })
      .then(({ data, errors }) => {
        console.log("Queued Job!");
      })
      .catch(err => console.log("Failed to queue job!", err));
  };

  return (
    <Card fluid>
      <Card.Content>
        <Card.Header as={Link} to={"/jobs/" + j.id}>
          {j.name}
        </Card.Header>
        <XTaskSummary
          tasks={j.tasks}
          limit={5}
          display={XTaskSummaryDisplayType.TARGET}
        />
        {j.staged && (
          <Button
            basic
            animated
            color="blue"
            size="small"
            style={{ margin: "15px" }}
            onClick={handleSubmit}
          >
            <Button.Content visible>{"Queue Job"}</Button.Content>
            <Button.Content hidden>
              <Icon name="arrow right" />
            </Button.Content>
          </Button>
        )}
      </Card.Content>
      <Card.Content extra>
        <XTags tags={j.tags} defaultText="None" />
      </Card.Content>
    </Card>
  );
};

export default XJobCard;
