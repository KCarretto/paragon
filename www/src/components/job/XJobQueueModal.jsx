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
    console.log("RENDERING MODAL")
    let container;

    const [isOpen, setIsOpen] = useState(false);
    const closeModal = () => {
        setIsOpen(false);
    }
    const openModal = () => {
        setIsOpen(true);
    }

    const [name, setName] = useState('');
    const [content, setContent] = useState('');
    const [tags, setTags] = useState([]);
    const [targets, setTargets] = useState([]);

    // const [params, setParams] = useState({ name: '', content: '', tags: [], targets: [] });
    // const handleChange = (e, { name, value }) => {
    //     let newParams = { [name]: value };
    //     console.log("Updated form values: ", name, value, params, newParams);
    //     setParams(newParams);
    // }

    const [queueJob, { called, loading, error }] = useMutation(QUEUE_JOB_MUTATION, {
        refetchQueries: [{ query: MULTI_JOB_QUERY }, { query: MULTI_TARGET_QUERY }],
    });

    const handleSubmit = () => {
        let vars = {
            name: name,
            content: content,
            tags: tags,
            targets: targets,
        }
        console.log("Creating job with vars: ", vars);

        queueJob({
            variables: vars
        }).then(({ data, errors }) => {
            // TODO: Display success or failure in form
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
                    value={name}
                    onChange={(e, { value }) => setName(value)}
                />
                <Form.Group widths={2}>
                    <Form.Field width={6}>
                        <label>Targets</label>
                        <XTargetTypeahead
                            onChange={(e, { value }) => setTargets(value)}
                        />
                    </Form.Field>
                    <Form.Field width={6}>
                        <label>Tags</label>
                        <XTagTypeahead
                            onChange={(e, { value }) => setTags(value)}
                        />
                    </Form.Field>
                </Form.Group>

                <Header inverted attached='top' size='large'>
                    <Icon name='code' />
                    {name ? name : 'Script'}
                </Header>
                <Form.Field control={XScriptEditor} onChange={(e, { value }) => setContent(value)} />

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
