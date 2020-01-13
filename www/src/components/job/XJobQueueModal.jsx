import { useMutation } from '@apollo/react-hooks';
import gql from 'graphql-tag';
import React, { useState } from 'react';
import { Button, Form, Header, Icon, Loader, Modal } from 'semantic-ui-react';
import { MULTI_JOB_QUERY, MULTI_TARGET_QUERY } from '../../views';
import { XScriptEditor, XTagTypeahead, XTargetTypeahead } from '../form';
import { XErrorMessage } from '../messages';

export const QUEUE_JOB_MUTATION = gql`
mutation QueueJob($name: String!, $content: String!, $tags: [ID!], $targets: [ID!]) {
    createJob(input: {name: $name, content: $content, tags: $tags, targets: $targets } ) {
        id
    }
}`;

const XJobQueueModal = ({ header }) => {
    let container;

    const [isOpen, setIsOpen] = useState(false);
    const closeModal = () => {
        setIsOpen(false);
    }
    const openModal = () => {
        setIsOpen(true);
    }

    const [params, setParams] = useState({ name: '', content: '', tags: [], targets: [] });
    const handleChange = (e, { name, value }) => {
        console.log("Updated form values: ", name, value);
        setParams({ ...params, [name]: value });
    }

    const [queueJob, { called, loading, error }] = useMutation(QUEUE_JOB_MUTATION, {
        refetchQueries: [{ query: MULTI_JOB_QUERY }, { query: MULTI_TARGET_QUERY }],
    });

    const handleSubmit = () => {
        console.log("Creating job with params: ", params);
        queueJob({ variables: params }).then(({ data, errors }) => {
            console.log("Create job result: ", data, errors);
            if (errors && errors.length > 0) {
                container.error(errors.join('\n', 'Failed to queue job'))
                return;
            }
            container.info(`Created job with id: ${data.job.id}`, 'Job Queued');
        }).catch((err) => console.error("GraphQL mutation failed", err));
    }

    console.log("ERROR IS", error);
    return (
        <Modal
            centered={false}
            open={isOpen}
            onClose={closeModal}
            trigger={<Button positive circular icon='plus' onClick={openModal} />}

            // Form properties
            as={Form}
            onSubmit={handleSubmit}
        >
            <Modal.Header>{header ? header : "Queue a Job"}<Loader disabled={!called || !loading} /></Modal.Header>
            <Modal.Content>
                <Form.Input
                    width={6}
                    label='Job Name'
                    placeholder='Enter job name'
                    name='name'
                    value={params.name}
                    onChange={handleChange}
                />
                <Form.Group widths={2}>
                    <Form.Field width={6}>
                        <label>Targets</label>
                        <XTargetTypeahead
                            onChange={handleChange}
                        />
                    </Form.Field>
                    <Form.Field width={6}>
                        <label>Tags</label>
                        <XTagTypeahead
                            onChange={handleChange}
                        />
                    </Form.Field>
                </Form.Group>

                <Header inverted attached='top' size='large'>
                    <Icon name='code' />
                    {params.name ? params.name : 'Script'}
                </Header>
                <XScriptEditor
                    handleChange={handleChange}
                />

                {/* <ControlledEditor
                    options={{
                        scrollbar: {
                            verticalScrollbarSize: '7px',
                        },
                        minimap: { enabled: false },
                        cursorStyle: 'line-thin',
                    }}
                    theme='dark'
                    height='600px'
                    value={params.content}
                    editorDidMount={(fn, mco) => {
                        let element = document.getElementsByTagName('textarea')[0];
                        element.classList.remove("inputarea");
                    }}
                    onChange={(e, value) => { handleChange(e, { name: 'content', value: value }) }}
                    language="python"
                /> */}

                {/* <Form.Field style={{ 'margin-top': '25px' }}>
                    <label>Script</label>
                    <Form.TextArea
                        label={{ content: 'Enter script' }}
                        placeholder='Enter script content'
                        name='content'
                        rows={15}
                        value={params.content}
                        onChange={handleChange}
                    />
                </Form.Field> */}
                <XErrorMessage title='Failed to Queue Job' err={error} />
            </Modal.Content>
            <Modal.Actions>
                <Form.Button positive floated='right'>Queue</Form.Button>
                <br />
            </Modal.Actions>
        </Modal >
    );
}

export default XJobQueueModal
