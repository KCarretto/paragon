import { useMutation, useQuery } from "@apollo/react-hooks";
import gql from "graphql-tag";
import * as React from "react";
import { useState } from "react";
import { DropdownItemProps, Form, Loader } from "semantic-ui-react";
import { Target } from "../../graphql/models";
import { MULTI_TAG_QUERY, MULTI_TARGET_QUERY } from "../../views";
import { MultiTagResult } from "../../views/XMultiTagView";
import { XErrorMessage } from "../messages";

const CREATE_TARGET = gql`
  mutation CreateTarget($name: String!, $primaryIP: String!, $tags: [ID!]) {
    createTarget(input: { name: $name, primaryIP: $primaryIP, tags: $tags }) {
      id
      name
      primaryIP
      lastSeen
      tags {
        id
        name
      }
    }
  }
`;

type Params = {
  name: String;
  primaryIP: String;
  tags: number[];
};

const XTargetCreateForm = () => {
  const [tagOptions, setTagOptions] = useState<DropdownItemProps[]>([]);

  const tagQueryResult = useQuery<MultiTagResult>(MULTI_TAG_QUERY, {
    onCompleted: data => {
      if (!data || !data.tags || data.tags.length < 1) {
        return;
      }
      setTagOptions(
        data.tags.map(({ id, name }) => {
          return {
            text: name,
            value: id
          };
        })
      );
      console.log("Updated tag suggestions: ", tagOptions);
    }
  });

  const [createTarget, { loading, error }] = useMutation<Target>(
    CREATE_TARGET,
    {
      refetchQueries: [{ query: MULTI_TARGET_QUERY }]
    }
  );

  const [params, setParams] = useState<Params>({
    name: "",
    primaryIP: "",
    tags: []
  });
  const handleChange = (e, { name, value }) => {
    setParams({ ...params, [name]: value });
  };
  const handleSubmit = () => {
    console.log("Creating target with params: ", params);
    createTarget({ variables: params }).then(({ data }) => {
      console.log("Create target result: ", data);
    });
  };

  return (
    <Form onSubmit={handleSubmit}>
      <Form.Input
        label={{ content: "Name" }}
        labelPosition="right"
        placeholder="Enter target name"
        name="name"
        value={params.name}
        onChange={handleChange}
      />
      <Form.Input
        label={{ content: "Primary IP" }}
        labelPosition="right"
        placeholder="Enter ip address"
        name="primaryIP"
        value={params.primaryIP}
        onChange={handleChange}
      />
      <Form.Dropdown
        placeholder="Add tags"
        fluid
        multiple
        search
        selection
        lazyLoad
        error={tagQueryResult.error}
        loading={tagQueryResult.loading}
        options={tagOptions}
        name="tags"
        value={params.tags}
        onChange={handleChange}
      />

      <Form.Button positive floated="right">
        Create
      </Form.Button>

      {error ? (
        <XErrorMessage title="Failed to create target" err={error} />
      ) : (
        <Loader active={loading} />
      )}
    </Form>
  );
};

export default XTargetCreateForm;
