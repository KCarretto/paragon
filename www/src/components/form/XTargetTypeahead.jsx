import { useQuery } from '@apollo/react-hooks';
import gql from 'graphql-tag';
import React, { useState } from 'react';
import { Dropdown, Input } from 'semantic-ui-react';

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
}`

// XTargetTypeahead adds a targets field to a form, which is an array of target ids. It provides
// tag suggestions as well, to allow the user to easily specify a set of targets with a tag.
// NOTE: Assumes global unique ids, no conflicts from tag & target ids.
const XTargetTypeahead = ({ onChange, labeled }) => {
    // optMap: Map of id => { text: name, value: id }
    // tagMap: Map of tag => [targets]
    // values: [tag id | target id]
    const [state, setState] = useState({ tagMap: new Map(), optMap: new Map(), values: [] });
    const [error, setError] = useState(null);

    // Wrap onChange to flatten target id array.
    const handleChange = (e, { name, value }) => {
        let targets = [...new Set(value.flatMap(id => {
            if (!state.tagMap.has(id)) {
                return id;
            }
            return state.tagMap.get(id);
        }))];

        setState({ ...state, values: value });
        onChange(e, { name: name, value: targets })
    }

    const { loading } = useQuery(SUGGEST_TARGETS_QUERY, {
        onCompleted: data => {
            if (!data || !data.targets) {
                data = { targets: [] };
            }

            let tMap = new Map();

            // For each target, create an Option for it and then create an additional option for
            // each of it's tags. Duplicate tags & targets will later be removed when the array of
            // [id, Option] is converted to a Map.
            let entries = data.targets.flatMap(target => [
                [
                    target.id,
                    {
                        text: target.name,
                        value: target.id
                    },
                ],
                ... !target.tags ? [] : target.tags.map(tag => {

                    // Maintain a map of tag => [targets]
                    if (!tMap.has(tag.id)) {
                        tMap.set(tag.id, [])
                    }
                    tMap.get(tag.id).push(target.id);

                    return [
                        tag.id,
                        {
                            text: tag.name,
                            value: tag.id,
                        }
                    ];
                })
            ]);

            setState({ ...state, tagMap: tMap, optMap: new Map(entries) });
            setError(null);
        },
        onError: err => {
            setError(err)
        }
    });

    let options = Array.from(state.optMap.values());
    const getDropdown = () => (
        <Dropdown
            placeholder='Select targets'
            icon=''
            fluid
            multiple
            search
            selection
            error={error}
            loading={loading}
            options={options}
            name='targets'
            value={state.values}
            onChange={handleChange}
            style={{
                borderRadius: "0 4px 4px 0",
            }}
        />
    );

    if (labeled) {
        return <Input
            fluid
            icon='desktop'
            label='Targets'
            input={getDropdown()}
        />
    }
    return getDropdown();
}

export default XTargetTypeahead;
