import { useMutation, useQuery } from "@apollo/react-hooks";
import gql from "graphql-tag";
import * as React from "react";
import { FunctionComponent, useState } from "react";
import { Link } from "react-router-dom";
import { Button, Icon, Input, Label, Table } from "semantic-ui-react";
import { SUGGEST_TAGS_QUERY, SUGGEST_TARGETS_QUERY, XTargetTypeahead } from "../components/form";
import { XBoundary } from "../components/layout";
import { XErrorMessage, XLoadingMessage } from "../components/messages";
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

  const [applyTag, { loading, error }] = useMutation(APPLY_TAG_TO_TARGET, {
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
    });
  };

  return (
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
  );
};

type TargetLabelProps = {
  tag: Tag;
  target: Target;
};

const TargetLabel: FunctionComponent<TargetLabelProps> = ({ tag, target }) => {
  const [removeTag] = useMutation(REMOVE_TAG_FROM_TARGET, {
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

    removeTag({
      variables: vars
    });
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

  const [createTag, { loading, error }] = useMutation(CREATE_TAG, {
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
    });
  };

  return (
    <Input
      fluid
      loading={loading}
      label={<Button icon="plus" positive onClick={handleSubmit} />}
      error={!!error}
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
  const whenEmpty = <span>No Targets</span>;

  return (
    <Table.Row>
      <Table.Cell width={2}>
        <Icon name="tag" />
        {tag.name}
      </Table.Cell>
      <Table.Cell collapsing width={8}>
        <XBoundary
          boundary={whenEmpty}
          show={tag && tag.targets && tag.targets.length > 0}
        >
          {tag && tag.targets && (
            <Label.Group style={{ maxWidth: "55vw", overflowX: "auto" }}>
              {tag.targets.map((target, index) => (
                <TargetLabel key={index} tag={tag} target={target} />
              ))}
            </Label.Group>
          )}
        </XBoundary>
      </Table.Cell>
      <Table.Cell width={4} singleLine collapsing textAlign="left">
        <AddTagToTargetsForm tag={tag} />
      </Table.Cell>
    </Table.Row>
  );
};

const XMultiTagView = () => {
  const { loading, error, data: { tags = [] } = {} } = useQuery<
    MultiTagResponse
  >(MULTI_TAG_QUERY, {
    pollInterval: 5000
  });

  const whenLoading = (
    <XLoadingMessage title="Loading Tags" msg="Fetching tag info" />
  );
  const whenEmpty = (
    <Table.Row>
      <Table.HeaderCell colSpan="3">
        <XNoTagsFound />
      </Table.HeaderCell>
    </Table.Row>
  );

  return (
    <React.Fragment>
      <XErrorMessage title="Error Loading Tags" err={error} />

      <XBoundary boundary={whenLoading} show={!loading}>
        <Table celled style={{ overflow: "auto" }}>
          <Table.Header>
            <Table.Row>
              <Table.HeaderCell>Tag</Table.HeaderCell>
              <Table.HeaderCell colSpan="2">Targets</Table.HeaderCell>
            </Table.Row>
          </Table.Header>

          <Table.Body>
            <XBoundary boundary={whenEmpty} show={tags && tags.length > 0}>
              {tags &&
                tags.map((tag, index) => (
                  <XTagTableRow key={index} tag={tag} />
                ))}
            </XBoundary>
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
      </XBoundary>
    </React.Fragment>
  );
};

export default XMultiTagView;
