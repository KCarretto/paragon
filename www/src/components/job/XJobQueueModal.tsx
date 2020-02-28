import { useMutation } from "@apollo/react-hooks";
import { ApolloError } from "apollo-client/errors/ApolloError";
import gql from "graphql-tag";
import * as React from "react";
import { useState } from "react";
import {
  Button,
  Checkbox,
  Form,
  Grid,
  Header,
  Icon,
  Input,
  Message,
  Modal
} from "semantic-ui-react";
import { Target } from "../../graphql/models";
import { MULTI_JOB_QUERY, MULTI_TARGET_QUERY } from "../../views";
import {
  useModal,
  XEditor,
  XServiceTypeahead,
  XTargetTypeahead
} from "../form";

export const QUEUE_JOB_MUTATION = gql`
  mutation QueueJob(
    $name: String!
    $content: String!
    $tags: [ID!]
    $targets: [ID!]
    $stage: Boolean
  ) {
    createJob(
      input: {
        name: $name
        content: $content
        tags: $tags
        targets: $targets
        stage: $stage
      }
    ) {
      id
    }
  }
`;

type JobQueueModalParams = {
  header?: string;
  openOnStart?: boolean;
};

const XJobQueueModal = ({ header, openOnStart }: JobQueueModalParams) => {
  const [openModal, closeModal, isOpen] = useModal();
  const [error, setError] = useState<ApolloError>(null);

  // Form params
  const [name, setName] = useState<string>("");
  const [content, setContent] = useState<string>(
    '\n# Enter your script here!\ndef main():\n\tprint("Hello World")'
  );
  // const [tags, setTags] = useState<Tag[]>([]);
  const [serviceTag, setServiceTag] = useState(null);
  const [targets, setTargets] = useState<Target[]>([]);
  const [stage, setStage] = useState<boolean>(false);

  const [queueJob, { called, loading }] = useMutation(QUEUE_JOB_MUTATION, {
    refetchQueries: [{ query: MULTI_JOB_QUERY }, { query: MULTI_TARGET_QUERY }]
  });

  const handleSubmit = () => {
    let vars = {
      name: name,
      content: content,
      tags: serviceTag ? [serviceTag] : [],
      targets: targets,
      stage: stage
    };

    queueJob({
      variables: vars
    })
      .then(({ data, errors }) => {
        if (errors && errors.length > 0) {
          let s = errors.map(e => e.message);
          let e = new ApolloError({
            graphQLErrors: errors,
            errorMessage: s.join("\n")
          });
          setError(e);
          return;
        }
        closeModal();
      })
      .catch(err => setError(err));
  };

  if (openOnStart) {
    openModal();
  }

  return (
    <Modal
      open={isOpen}
      onClose={closeModal}
      trigger={<Button icon="cube" color="green" onClick={openModal} />}
      size="large"
      // Form properties
      as={Form}
      onSubmit={handleSubmit}
      error={called && error}
      loading={called && loading}
    >
      <Modal.Header>{header ? header : "Queue a Job"}</Modal.Header>
      <Modal.Content>
        <Grid verticalAlign="middle" stackable container columns={"equal"}>
          <Grid.Row>
            <Grid.Column>
              <Input
                label="Job Name"
                icon="cube"
                fluid
                placeholder="Enter job name"
                name="name"
                value={name}
                onChange={(e, { value }) => setName(value)}
              />
            </Grid.Column>

            <Grid.Column>
              <XServiceTypeahead
                labeled
                value={serviceTag}
                onChange={(e, { value }) => setServiceTag(value)}
              />
              {/* <XTagTypeahead
              labeled
              onChange={(e, { value }) => setTags(value)}
            /> */}
            </Grid.Column>

            <Grid.Column>
              <XTargetTypeahead
                labeled
                onChange={(e, { value }) => setTargets(value)}
              />
            </Grid.Column>
          </Grid.Row>
          <Grid.Row>
            <Grid.Column>
              <Checkbox
                label="Stage this job"
                onChange={() => setStage(!stage)}
                checked={stage}
              />
            </Grid.Column>
          </Grid.Row>
        </Grid>

        <Header inverted attached="top" size="large">
          <Icon name="code" />
          {name ? name : "Script"}
        </Header>
        <Form.Field
          control={XEditor}
          value={content}
          onChange={(e, value) => setContent(value)}
          // control={XScriptEditor}
          // content={content}
          // onChange={(e, { value }) => setContent(value)}
        />
        <Message
          error
          icon="warning"
          header={"Failed to Queue Job"}
          onDismiss={(e, data) => setError(null)}
          content={error ? error.message : "Unknown Error"}
        />
      </Modal.Content>
      <Modal.Actions>
        <Form.Button style={{ marginBottom: "10px" }} positive floated="right">
          {stage ? "Stage Job" : "Queue"}
        </Form.Button>
      </Modal.Actions>
    </Modal>
  );
};

export default XJobQueueModal;
