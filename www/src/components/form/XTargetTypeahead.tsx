import { useQuery } from "@apollo/react-hooks";
import { ApolloError } from "apollo-client/errors/ApolloError";
import gql from "graphql-tag";
import * as React from "react";
import { useState } from "react";
import { Dropdown, DropdownItemProps, Input } from "semantic-ui-react";
import { Tag, Target } from "../../graphql/models";

// Suggest targets and tags, but only suggest tags that have at least one target.
export const SUGGEST_TARGETS_QUERY = gql`
  query SuggestTargets {
    targets {
      id
      name
      tags {
        id
        name
      }
    }
  }
`;

type State = {
  optMap: Map<string, DropdownItemProps>;
  tagMap: Map<Tag, Target[]>;
  values: string[];
};

type TargetsResult = {
  targets: Target[];
};

// XTargetTypeahead adds a targets field to a form, which is an array of target ids. It provides
// tag suggestions as well, to allow the user to easily specify a set of targets with a tag.
// NOTE: Assumes global unique ids, no conflicts from tag & target ids.
const XTargetTypeahead = ({ onChange, labeled }) => {
  // optMap: Map of id => { text: name, value: id }
  // tagMap: Map of tag => [targets]
  // values: [tag id | target id]
  const [state, setState] = useState<State>({
    tagMap: new Map(),
    optMap: new Map(),
    values: []
  });
  const [error, setError] = useState<ApolloError>(null);

  // Wrap onChange to flatten target id array.
  const handleChange = (e, { name, value }) => {
    let targets = [
      ...new Set(
        value.flatMap(id => {
          if (!state.tagMap.has(id)) {
            return id;
          }
          return state.tagMap.get(id);
        })
      )
    ];

    setState({ ...state, values: value });
    onChange(e, { name: name, value: targets });
  };

  const { loading } = useQuery<TargetsResult>(SUGGEST_TARGETS_QUERY, {
    onCompleted: data => {
      if (!data || !data.targets) {
        data = { targets: [] };
      }

      let tMap = new Map();

      // For each target, create an Option for it and then create an additional option for
      // each of it's tags. Duplicate tags & targets will later be removed when the array of
      // [id, Option] is converted to a Map.

      let entries = [];

      data.targets.map((target): void => {
        entries.push([
          target.id,
          {
            text: target.name,
            value: target.id
          }
        ]);
        if (target.tags !== null) {
          target.tags.map((tag): void => {
            if (!tMap.has(tag.id)) {
              tMap.set(tag.id, []);
            }
            tMap.get(tag.id).push(target.id);
            entries.push([
              tag.id,
              {
                text: tag.name,
                value: tag.id
              }
            ]);
            return null;
          });
        }
        return null;
      });

      setState({ ...state, tagMap: tMap, optMap: new Map(entries) });
      setError(null);
    },
    onError: err => {
      setError(err);
    }
  });

  let options = Array.from(state.optMap.values());
  const getDropdown = () => (
    <Dropdown
      placeholder="Select targets"
      icon=""
      fluid
      multiple
      search
      selection
      error={error !== null}
      loading={loading}
      options={options}
      name="targets"
      value={state.values}
      onChange={handleChange}
      style={{
        borderRadius: "0 4px 4px 0"
      }}
    />
  );

  if (labeled) {
    return <Input fluid icon="desktop" label="Targets" input={getDropdown()} />;
  }
  return getDropdown();
};

export default XTargetTypeahead;
