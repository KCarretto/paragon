import { useMutation, useQuery } from '@apollo/react-hooks';
import gql from 'graphql-tag';
import React, { useState } from 'react';
import { Form, Loader } from 'semantic-ui-react';
import { MULTI_JOB_QUERY, MULTI_TAG_QUERY, MULTI_TARGET_QUERY } from '../../views';
import { XErrorMessage } from '../messages';

const QUEUE_JOB = gql`
mutation QueueJob($name: String!, $content: String!, $tags: [ID!]) {
    createJob(input: {name: $name, content: $content, tags: $tags } ) {
        id
    }
}`;

const XJobQueueForm = () => {
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

    const [queueJob, { loading, error }] = useMutation(QUEUE_JOB, {
        refetchQueries: [{ query: MULTI_JOB_QUERY }, { query: MULTI_TARGET_QUERY }],
    });

    const [params, setParams] = useState({ name: '', content: '', tags: [] });
    const handleChange = (e, { name, value }) => {
        setParams({ ...params, [name]: value });
    }
    const handleSubmit = () => {
        console.log("Creating job with params: ", params);
        queueJob({ variables: params }).then(({ data }) => {
            console.log("Create job result: ", data);
        });
    }

    return (
        <Form onSubmit={handleSubmit}>
            <Form.Input
                label={{ content: 'Name' }}
                labelPosition='right'
                placeholder='Enter job name'
                name='name'
                value={params.name}
                onChange={handleChange}
            />
            <Form.TextArea
                label={{ content: 'Enter script' }}
                labelPosition='right'
                placeholder='Enter script content'
                name='content'
                rows={15}
                value={params.content}
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

            <Form.Button positive floated='right'>Queue</Form.Button>

            {error ?
                <XErrorMessage title='Failed to queue job' msg={`${error}`} />
                : <Loader active={loading} />}
        </Form>
    );
}

export default XJobQueueForm;
