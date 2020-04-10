import { useMutation } from "@apollo/react-hooks";
import gql from "graphql-tag";
import * as React from "react";
import { useState } from "react";
import { toast } from "react-semantic-toasts";
import { Button, Grid, Icon, Input, Segment } from "semantic-ui-react";
import { XJobSettingsModal } from ".";
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
    errorPolicy: "all",
    refetchQueries: [{ query: RUNS_QUERY, variables: { name: name } }]
  });

  const onSubmit = (vars: CreateJobRequest) => {
    let valid = true;
    if (!targets || targets.length < 1) {
      toast({
        title: "凸(¬‿¬)凸 Must Select Targets",
        type: "error",
        description: "Choose some targets before running the job",
        time: 5000
      });
      valid = false;
    }
    if (!name || name === "") {
      toast({
        title: "凸(¬‿¬)凸 No Job Name Provided",
        type: "error",
        description: "Give the job a name before running it",
        time: 5000
      });
      valid = false;
    }
    if (!content || content === "") {
      toast({
        title: "凸(¬‿¬)凸 Empty Script Provided",
        type: "error",
        description: "We imagine you actually wanted to run something...",
        time: 5000
      });
      valid = false;
    }
    if (!valid) {
      return;
    }

    createJob({
      variables: vars
    })
      .then(() => {
        toast({
          title: "(◠ω◠✿) Job Queued",
          type: "success",
          description: "Check below for execution results",
          time: 3000
        });
      })
      .catch(err => {
        console.log("[ERROR] FAILED TO QUEUE JOB", err);
        toast({
          title: "(҂◡_◡) Failed to Queue Job",
          type: "error",
          description: `${err}`,
          time: 5000
        });
      });
  };

  return (
    <React.Fragment>
      <Segment basic secondary inverted attached="top" style={{ borderRadius: "0px" }}>
        <Grid stackable verticalAlign="middle">
          {/* Job Name */}
          <Grid.Column width={6}>
            <Input
              icon="code"
              iconPosition="left"
              placeholder="Name this job..."
              transparent
              inverted
              size="huge"
              name="name"
              value={name}
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
                loading={loading}
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
              <XJobSettingsModal
                serviceTag={serviceTag}
                setServiceTag={setServiceTag}
                stage={stage}
                setStage={setStage}
              />
              {/* <Button icon="cogs" color="teal" /> */}
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
