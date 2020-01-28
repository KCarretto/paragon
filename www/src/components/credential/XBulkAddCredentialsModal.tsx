import { useMutation } from "@apollo/react-hooks";
import { ApolloError } from "apollo-client/errors/ApolloError";
import gql from "graphql-tag";
import * as React from "react";
import { useState } from "react";
import {
  Button,
  Form,
  Grid,
  Header,
  Icon,
  Input,
  Message,
  Modal
} from "semantic-ui-react";
import { Tag, Target } from "../../graphql/models";
import { MULTI_TARGET_QUERY } from "../../views";
import {
  useModal,
  XScriptEditor,
  XTagTypeahead,
  XTargetTypeahead
} from "../form";

export const QUEUE_JOB_MUTATION = gql`
  mutation QueueJob(
    $name: String!
    $content: String!
    $tags: [ID!]
    $targets: [ID!]
  ) {
    createJob(
      input: { name: $name, content: $content, tags: $tags, targets: $targets }
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
  const [tags, setTags] = useState<Tag[]>([]);
  const [targets, setTargets] = useState<Target[]>([]);

  const [queueJob, { called, loading }] = useMutation(QUEUE_JOB_MUTATION, {
    refetchQueries: [{ query: MULTI_TARGET_QUERY }]
  });

  const handleSubmit = () => {
    let vars = {
      name: name,
      content: content,
      tags: tags,
      targets: targets
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
      trigger={<Button positive circular icon="plus" onClick={openModal} />}
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
            <XTagTypeahead
              labeled
              onChange={(e, { value }) => setTags(value)}
            />
          </Grid.Column>

          <Grid.Column>
            <XTargetTypeahead
              labeled
              onChange={(e, { value }) => setTargets(value)}
            />
          </Grid.Column>
        </Grid>

        <Header inverted attached="top" size="large">
          <Icon name="code" />
          {name ? name : "Script"}
        </Header>
        <Form.Field
          control={XScriptEditor}
          content={content}
          onChange={(e, { value }) => setContent(value)}
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
          Queue
        </Form.Button>
      </Modal.Actions>
    </Modal>
  );
};

export default XJobQueueModal;
