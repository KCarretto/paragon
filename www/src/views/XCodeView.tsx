import { useMutation, useQuery } from "@apollo/react-hooks";
import { ApolloError } from "apollo-client/errors/ApolloError";
import gql from "graphql-tag";
import * as React from "react";
import { useState } from "react";
import { useParams } from "react-router-dom";
import {
  Button,
  Form,
  Grid,
  Header,
  Icon,
  Input,
  Segment
} from "semantic-ui-react";
import { MULTI_JOB_QUERY, MULTI_TARGET_QUERY } from ".";
import { XEditor, XTargetTypeahead } from "../components/form";
import { XBoundary, XCardGroup } from "../components/layout";
import { XErrorMessage, XLoadingMessage } from "../components/messages";
import { XTaskCard, XTaskCardDisplayType } from "../components/task";
import { CreateJobRequest, Job, Tag } from "../graphql/models";

export const JOB_QUERY = gql`
  query Job($id: ID!) {
    job(id: $id) {
      id
      name
      content
      staged
      tags {
        id
        name
      }
      tasks {
        id
        queueTime
        claimTime
        execStartTime
        execStopTime
        content
        output
        error
        sessionID

        target {
          id
          name
        }
        job {
          id
          name
          staged
          tags {
            id
            name
          }
        }
      }
    }
  }
`;

type JobQuery = {
  job: Job;
};

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

const XJobEditor: React.FC<{
  prevName: string;
  prevContent: string;
  prevTags: Tag[];
  prevStaged: boolean;
  onSubmit: (vars: CreateJobRequest) => void;
}> = ({ prevName, prevContent, prevTags, prevStaged, onSubmit }) => {
  const [name, setName] = useState<string>(prevName);
  const [content, setContent] = useState<string>(
    prevContent ||
      '\n# Enter your script here!\ndef main():\n\tprint("Hello World")'
  );
  const [serviceTag, setServiceTag] = useState(
    prevTags && prevTags.length > 0 ? prevTags[0].id : null
  );
  const [targets, setTargets] = useState<string[]>([]);
  const [stage, setStage] = useState<boolean>(prevStaged);

  return (
    <React.Fragment>
      {/* Code Editor Header */}
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
      <Form.Field
        control={XEditor}
        value={content}
        onChange={(e, value) => setContent(value)}
      />
    </React.Fragment>
  );
};

const XJobResults: React.FC<{ id: string }> = ({ id }) => {
  const { loading, error, data: { job: { tasks = [] } = {} } = {} } = useQuery<
    JobQuery
  >(JOB_QUERY, {
    skip: !id,
    variables: { id: id },
    pollInterval: 2500
  });

  return (
    <React.Fragment>
      <XBoundary
        show={!loading}
        boundary={
          <XLoadingMessage
            title="Loading Results"
            msg="Fetching execution results"
          />
        }
      >
        <XBoundary
          show={!error}
          boundary={
            <XErrorMessage title="Failed to Load Results" err={error} />
          }
        >
          {tasks && tasks.length > 0 ? (
            <XCardGroup>
              {tasks.map(task => (
                <XTaskCard
                  key={task.id}
                  display={XTaskCardDisplayType.TARGET}
                  task={task}
                />
              ))}
            </XCardGroup>
          ) : (
            <p>Run your job to see output for each target system here.</p>
          )}
        </XBoundary>
      </XBoundary>
    </React.Fragment>
  );
};

const XCodeView = () => {
  let { id: idQueryParam = null } = useParams();

  const {
    loading,
    error,
    data: {
      job: { name = null, content = null, tags = [], staged = false } = {}
    } = {}
  } = useQuery<JobQuery>(JOB_QUERY, {
    skip: !idQueryParam,
    variables: { id: idQueryParam },
    pollInterval: 5000
  });

  const [jobID, setJobID] = useState<string>(idQueryParam);

  const [queueJob, { called, loading: mutationLoading }] = useMutation(
    QUEUE_JOB_MUTATION,
    {
      refetchQueries: [
        { query: MULTI_JOB_QUERY },
        { query: MULTI_TARGET_QUERY }
      ]
    }
  );
  console.log(mutationLoading);

  const onSubmit = (vars: CreateJobRequest) => {
    queueJob({
      variables: vars
    })
      .then(({ data, errors }) => {
        setJobID(data && data.job && data.job.id ? data.job.id : jobID);

        if (errors && errors.length > 0) {
          let s = errors.map(e => e.message);
          let e = new ApolloError({
            graphQLErrors: errors,
            errorMessage: s.join("\n")
          });
          alert(e);
          return;
        }
      })
      .catch(err => alert(err));
  };

  return (
    <React.Fragment>
      <XBoundary
        show={!idQueryParam || !!jobID || !loading} // Show if no id queryParam, we have the prevJobID already, or we're finished loading
        boundary={
          <XLoadingMessage
            title="Loading Editor"
            msg="Fetching previous run information"
          />
        }
      >
        <XJobEditor
          prevName={name}
          prevContent={content}
          prevTags={tags}
          prevStaged={staged}
          onSubmit={onSubmit}
        />
        <Header size="large" block inverted>
          <Icon name="tasks" />
          <Header.Content>Results</Header.Content>
        </Header>
        <XJobResults id={jobID} /> {/* TODO: Needs to query by name */}
      </XBoundary>
    </React.Fragment>
  );
};

export default XCodeView;
