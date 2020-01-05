import { useMutation, useQuery } from '@apollo/react-hooks';
import gql from 'graphql-tag';
import React, { useState } from 'react';
import { Form, Loader } from 'semantic-ui-react';
import { MULTI_TAG_QUERY, MULTI_TARGET_QUERY } from '../../views';
import { XErrorMessage } from '../messages';

const CREATE_TARGET = gql`
mutation CreateTarget($name: String!, $primaryIP: String!, $tags: [ID!]) {
    createTarget(input: {name: $name, primaryIP: $primaryIP, tags: $tags } ) {
        id
        name
        primaryIP
        lastSeen
        tags {
            id
            name
        }
    }
}`;

const XTargetCreateForm = () => {
    const [tagOptions, setTagOptions] = useState([]);

    const { tagLoading, tagError } = useQuery(MULTI_TAG_QUERY, {
        onCompleted: data => {
            if (!data || !data.tags || data.tags.length < 1) {
                return;
            }
            setTagOptions(data.tags.map(({ id, name }) => {
                return {
                    text: name,
                    value: id,
                };
            }));
            console.log("Updated tag suggestions: ", tagOptions)
        },
    });

    const [createTarget, { loading, error }] = useMutation(CREATE_TARGET, {
        refetchQueries: [{ query: MULTI_TARGET_QUERY }],
    });

    const [params, setParams] = useState({ name: '', primaryIP: '', tags: [] });
    const handleChange = (e, { name, value }) => {
        setParams({ ...params, [name]: value });
    }
    const handleSubmit = () => {
        console.log("Creating target with params: ", params);
        createTarget({ variables: params }).then(({ data }) => {
            console.log("Create target result: ", data);
        });
    }



    return (
        <Form onSubmit={handleSubmit}>
            <Form.Input
                label={{ content: 'Name' }}
                labelPosition='right'
                placeholder='Enter target name'
                name='name'
                value={params.name}
                onChange={handleChange}
            />
            <Form.Input
                label={{ content: 'Primary IP' }}
                labelPosition='right'
                placeholder='Enter ip address'
                name='primaryIP'
                value={params.primaryIP}
                onChange={handleChange}
            />
            <Form.Dropdown
                placeholder='Add tags'
                fluid
                multiple
                search
                selection
                lazyLoad
                error={tagError}
                loading={tagLoading}
                options={tagOptions}
                name='tags'
                value={params.tags}
                onChange={handleChange}
            />

            <Form.Button positive floated='right'>Create</Form.Button>

            {error ?
                <XErrorMessage title='Failed to create target' msg={`${error}`} />
                : <Loader active={loading} />}
        </Form>
    );
}

export default XTargetCreateForm;
