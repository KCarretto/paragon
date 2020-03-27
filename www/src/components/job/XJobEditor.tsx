import { useMutation } from "@apollo/react-hooks";
import gql from "graphql-tag";
import * as React from "react";
import { useState } from "react";
import { Button, Grid, Icon, Input, Segment } from "semantic-ui-react";
import { CreateJobRequest } from "../../graphql/models";
import { XEditor, XTargetTypeahead } from "../form";
import { RUNS_QUERY } from "./XJobResults";

export const CREATE_JOB_MUTATION = gql`
  mutation CreateJob(
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

const XJobEditor: React.FC<{ name: string; setName: (string) => void }> = ({
  name,
  setName
}) => {
  //   const [name, setName] = useState<string>("Untitled Job");
  const [content, setContent] = useState<string>(
    '\n# Enter your script here!\ndef main():\n\tprint("Hello World")'
  );
  const [serviceTag, setServiceTag] = useState<string>(null);
  const [targets, setTargets] = useState<string[]>([]);
  const [stage, setStage] = useState<boolean>(false);

  const [createJob, { loading }] = useMutation(CREATE_JOB_MUTATION, {
    refetchQueries: [{ query: RUNS_QUERY, variables: { name: name } }]
  });

  const onSubmit = (vars: CreateJobRequest) => {
    createJob({
      variables: vars
    });
  };

  return (
    <React.Fragment>
      <Segment basic secondary inverted attached="top">
        <Grid stackable verticalAlign="middle">
          {/* Job Name */}
          <Grid.Column width={6}>
            <Input
              icon="code"
              iconPosition="left"
              transparent
              inverted
              size="huge"
              name="name"
              value={name || "Untitled Job..."}
              onChange={(e, { value }) => setName(value)}
            />
          </Grid.Column>

          {/* Target Selection */}
          <Grid.Column width={8} textAlign="right" floated="right">
            <XTargetTypeahead
              labeled
              onChange={(e, { value }) => setTargets(value)}
            />
          </Grid.Column>

          {/* Action Buttons */}
          <Grid.Column width={2} textAlign="right" floated="right">
            <Button.Group>
              <Button
                animated
                color="blue"
                onClick={e =>
                  onSubmit({
                    name: name,
                    content: content,
                    targets: targets,
                    tags: serviceTag ? [serviceTag] : [],
                    stage: stage
                  })
                }
              >
                <Button.Content visible>Run</Button.Content>
                <Button.Content hidden>
                  <Icon name="arrow right" />
                </Button.Content>
              </Button>
              <Button icon="cogs" color="teal" />
            </Button.Group>
          </Grid.Column>
        </Grid>
      </Segment>

      {/* Code Editor */}
      <XEditor value={content} onChange={(e, value) => setContent(value)} />
    </React.Fragment>
  );
};

export default XJobEditor;
