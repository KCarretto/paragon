import { useMutation, useQuery } from "@apollo/react-hooks";
import { ApolloError } from "apollo-client/errors/ApolloError";
import gql from "graphql-tag";
import * as React from "react";
import { FunctionComponent, useState } from "react";
import { Link } from "react-router-dom";
import { Button, Icon, Input, Loader, Table } from "semantic-ui-react";
import {
  SUGGEST_TAGS_QUERY,
  SUGGEST_TARGETS_QUERY,
  XTargetTypeahead
} from "../components/form";
import { XErrorMessage } from "../components/messages";
import { XNoTagsFound } from "../components/tag";
import { Tag, Target } from "../graphql/models";

export const MULTI_TAG_QUERY = gql`
  {
    tags {
      id
      name
      targets {
        id
        name
      }
    }
  }
`;

const CREATE_TAG = gql`
  mutation CreateTag($name: String!) {
    createTag(input: { name: $name }) {
      id
      name
    }
  }
`;

const APPLY_TAG_TO_TARGET = gql`
  mutation ApplyTag($tag: ID!, $targets: [ID!]) {
    applyTagToTargets(input: { tagID: $tag, targets: $targets }) {
      id
      name
    }
  }
`;

const REMOVE_TAG_FROM_TARGET = gql`
  mutation RemoveTag($tag: ID!, $target: ID!) {
    removeTagFromTarget(input: { tagID: $tag, entID: $target }) {
      id
      name
    }
  }
`;

// export const ADD_TAG_TO_TARGETS_MUTATION = gql`
// {

// }`;

export type MultiTagResponse = {
  tags: Tag[];
};

type AddTagToTargetsFormProps = {
  tag: Tag;
};

const AddTagToTargetsForm: FunctionComponent<AddTagToTargetsFormProps> = ({
  tag
}) => {
  const [targets, setTargets] = useState<Target[]>([]);
  const [error, setError] = useState<ApolloError>(null);

  const [applyTag, { called, loading }] = useMutation(APPLY_TAG_TO_TARGET, {
    refetchQueries: [
      { query: MULTI_TAG_QUERY },
      { query: SUGGEST_TAGS_QUERY },
      { query: SUGGEST_TARGETS_QUERY }
    ]
  });

  const handleSubmit = () => {
    let vars = {
      tag: tag.id,
      targets: targets
    };

    applyTag({
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
      })
      .catch(err => setError(err));
  };

  return (
    // <Input
    //   fluid
    //   label={}
    //   labelPosition="right"
    // input={
    <XTargetTypeahead
      onChange={(e, { value }) => setTargets(value)}
      labeled
      input={{
        loading: loading,
        error: error !== null,
        fluid: true,
        label: <Button icon="plus" positive onClick={handleSubmit} />,
        labelPosition: "right"
      }}
    />
    // }
    // />
  );
};

type TargetLabelProps = {
  tag: Tag;
  target: Target;
};

const TargetLabel: FunctionComponent<TargetLabelProps> = ({ tag, target }) => {
  const [error, setError] = useState<ApolloError>(null);

  const [removeTag, { called, loading }] = useMutation(REMOVE_TAG_FROM_TARGET, {
    refetchQueries: [
      { query: MULTI_TAG_QUERY },
      { query: SUGGEST_TAGS_QUERY },
      { query: SUGGEST_TARGETS_QUERY }
    ]
  });

  const handleRemove = () => {
    let vars = {
      tag: tag.id,
      target: target.id
    };
    console.log("APPLY VARS", vars);

    removeTag({
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
      })
      .catch(err => setError(err));
  };
  return (
    <Button.Group labeled>
      <Button compact as={Link} to={"/targets/" + target.id}>
        {target.name}
      </Button>
      <Button compact negative icon="x" onClick={handleRemove} />
    </Button.Group>
  );
};

const AddTagForm = () => {
  const [name, setName] = useState("");
  const [error, setError] = useState<ApolloError>(null);

  const [createTag, { called, loading }] = useMutation(CREATE_TAG, {
    refetchQueries: [
      { query: MULTI_TAG_QUERY },
      { query: SUGGEST_TAGS_QUERY },
      { query: SUGGEST_TARGETS_QUERY }
    ]
  });

  const handleSubmit = () => {
    let vars = {
      name: name
    };

    createTag({
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
        setName("");
      })
      .catch(err => setError(err));
  };

  return (
    <Input
      fluid
      loading={loading}
      label={<Button icon="plus" positive onClick={handleSubmit} />}
      error={error !== null}
      labelPosition="right"
      iconPosition="left"
      icon="tag"
      placeholder="Create a tag..."
      name="name"
      value={name}
      onChange={(e, { value }) => setName(value)}
    />
  );
};

type XTagTableRowProps = {
  tag: Tag;
};

const XTagTableRow: FunctionComponent<XTagTableRowProps> = ({ tag }) => {
  return (
    <Table.Row>
      <Table.Cell width={2}>
        <Icon name="tag" />
        {tag.name}
      </Table.Cell>
      <Table.Cell collapsing width={8}>
        {tag.targets && tag.targets.length > 0
          ? tag.targets.map((target, index) => (
              <TargetLabel key={index} tag={tag} target={target} />
            ))
          : "No Targets"}
      </Table.Cell>
      <Table.Cell width={4} singleLine collapsing textAlign="left">
        <AddTagToTargetsForm tag={tag} />
      </Table.Cell>
    </Table.Row>
  );
};

const XMultiTagView = () => {
  const { called, loading, error, data } = useQuery<MultiTagResponse>(
    MULTI_TAG_QUERY
  );

  const showList = () => {
    return (
      <Table celled>
        <Table.Header>
          <Table.Row>
            <Table.HeaderCell>Tag</Table.HeaderCell>
            <Table.HeaderCell colSpan="2">Targets</Table.HeaderCell>
          </Table.Row>
        </Table.Header>

        <Table.Body>
          {data && data.tags && data.tags.length > 0 ? (
            data.tags.map((tag, index) => (
              <XTagTableRow key={index} tag={tag} />
            ))
          ) : (
            <Table.Row>
              <Table.HeaderCell colSpan="3">
                <XNoTagsFound />
              </Table.HeaderCell>
            </Table.Row>
          )}
        </Table.Body>
        <Table.Footer fullWidth>
          <Table.Row>
            <Table.HeaderCell width={12} colSpan="2" />
            <Table.HeaderCell width={4} textAlign="right">
              <AddTagForm />
            </Table.HeaderCell>
          </Table.Row>
        </Table.Footer>
      </Table>
    );
  };

  return (
    <React.Fragment>
      <Loader disabled={!called || !loading} />

      {showList()}

      <XErrorMessage title="Error Loading Tags" err={error} />
    </React.Fragment>
  );
};

export default XMultiTagView;
