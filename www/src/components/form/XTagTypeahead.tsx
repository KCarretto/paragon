import { useQuery } from "@apollo/react-hooks";
import { ApolloError } from "apollo-client/errors/ApolloError";
import gql from "graphql-tag";
import * as React from "react";
import { useState } from "react";
import { Dropdown, DropdownItemProps, Input } from "semantic-ui-react";
import { Tag } from "../../graphql/models";

// Suggest tags for the typeahead.
export const SUGGEST_TAGS_QUERY = gql`
  query SuggestTags {
    tags {
      id
      name
    }
  }
`;

type TagsResult = {
  tags: Tag[];
};

type State = {
  optMap: Map<string, DropdownItemProps>;
  values: string[];
};

// XTagTypeahead adds a tags field to a form, which is an array of tag ids with no duplicates.
const XTagTypeahead = ({ onChange, labeled }) => {
  // Map of id => { text: name, value: id }
  const [state, setState] = useState<State>({
    optMap: new Map(),
    values: []
  });
  const [error, setError] = useState<ApolloError>(null);

  const handleChange = (e, { name, value }) => {
    setState({ ...state, values: value });
    onChange(e, { name: name, value: value });
  };

  // NOTE: Assumes global unique ids, no conflicts from tag & target ids.
  const { loading } = useQuery<TagsResult>(SUGGEST_TAGS_QUERY, {
    onCompleted: data => {
      if (!data || !data.tags) {
        data = { tags: [] };
      }

      // Build a map of tag.id => Option
      let tags: [string, DropdownItemProps][] = data.tags.map(tag => [
        tag.id,
        {
          text: tag.name,
          value: tag.id
        }
      ]);

      setState({ ...state, optMap: new Map(tags) });
      setError(null);
    },
    onError: err => {
      setError(err);
    }
  });

  let options = Array.from(state.optMap.values());
  const getDropdown = () => (
    <Dropdown
      placeholder="Add tags"
      icon=""
      fluid
      multiple
      search
      selection
      error={error !== null}
      loading={loading}
      options={options}
      name="tags"
      value={state.values}
      onChange={handleChange}
      style={{
        borderRadius: "0 4px 4px 0"
      }}
    />
  );

  if (labeled) {
    return <Input fluid label="Tags" icon="tags" input={getDropdown()} />;
  }
  return getDropdown();
};

export default XTagTypeahead;
