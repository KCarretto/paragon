import { useQuery } from '@apollo/react-hooks';
import gql from 'graphql-tag';
import React, { useState } from 'react';
import { Dropdown, Input } from 'semantic-ui-react';

// Suggest tags for the typeahead.
export const SUGGEST_TAGS_QUERY = gql`
query SuggestTags {
    tags {
        id
        name
    }
}`

// XTagTypeahead adds a tags field to a form, which is an array of tag ids with no duplicates.
const XTagTypeahead = ({ onChange, labeled }) => {
    // Map of id => { text: name, value: id }
    const [state, setState] = useState({ optMap: new Map(), values: [] });

    const handleChange = (e, { name, value }) => {
        setState({ ...state, values: value });
        onChange(e, { name: name, value: value })
    }

    // NOTE: Assumes global unique ids, no conflicts from tag & target ids.
    const { loading, err } = useQuery(SUGGEST_TAGS_QUERY, {
        onCompleted: data => {
            if (!data || !data.tags) {
                data = { tags: [] };
            }

            // Build a map of tag.id => Option
            let tags = data.tags.map(tag => [
                tag.id,
                {
                    text: tag.name,
                    value: tag.id
                }
            ]);

            setState({ ...state, optMap: new Map(tags) });
        }
    });

    let options = Array.from(state.optMap.values());
    const getDropdown = () => (
        <Dropdown
            placeholder='Add tags'
            icon=''
            fluid
            multiple
            search
            selection
            error={err}
            loading={loading}
            options={options}
            name='tags'
            value={state.values}
            onChange={handleChange}
            style={{
                borderRadius: "0 4px 4px 0",
            }}
        />
    );

    if (labeled) {
        return (
            <Input
                fluid
                label='Tags'
                icon='tags'
                input={getDropdown()}
            />
        );
    }
    return getDropdown();
}

export default XTagTypeahead;
