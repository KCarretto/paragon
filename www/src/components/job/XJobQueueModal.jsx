import { useMutation } from '@apollo/react-hooks';
import gql from 'graphql-tag';
import React, { useState } from 'react';
import { Button, Form, Grid, Header, Icon, Input, Message, Modal } from 'semantic-ui-react';
import { MULTI_JOB_QUERY, MULTI_TARGET_QUERY } from '../../views';
import { useModal, XScriptEditor, XTagTypeahead, XTargetTypeahead } from '../form';

export const QUEUE_JOB_MUTATION = gql`
mutation QueueJob($name: String!, $content: String!, $tags: [ID!], $targets: [ID!]) {
    createJob(input: {name: $name, content: $content, tags: $tags, targets: $targets } ) {
        id
    }
}`;

const XJobQueueModal = ({ header }) => {
    const [openModal, closeModal, isOpen] = useModal();
    const [err, setError] = useState(null);

    const [name, setName] = useState('');
    const [content, setContent] = useState('\n# Enter your script here!\ndef main():\n\tprint("Hello World")');
    const [tags, setTags] = useState([]);
    const [targets, setTargets] = useState([]);

    const [queueJob, { called, loading }] = useMutation(QUEUE_JOB_MUTATION, {
        refetchQueries: [{ query: MULTI_JOB_QUERY }, { query: MULTI_TARGET_QUERY }],
    });

    const handleSubmit = () => {
        let vars = {
            name: name,
            content: content,
            tags: tags,
            targets: targets,
        }

        queueJob({
            variables: vars
        }).then(({ data, errors }) => {
            if (errors && errors.length > 0) {
                setError({ message: errors.join('\n') })
                return;
            }
            closeModal();
        }).catch((error) => setError(error));
    }

    return (
        <Modal
            open={isOpen}
            onClose={closeModal}
            trigger={<Button positive circular icon='plus' onClick={openModal} />}
            size='large'

            // Form properties
            as={Form}
            onSubmit={handleSubmit}
            error={called && err}
            loading={called && loading}
        >
            <Modal.Header>{header ? header : "Queue a Job"}</Modal.Header>
            <Modal.Content>
                <Grid verticalAlign='middle' stackable container columns={'equal'}>
                    <Grid.Column>
                        <Input
                            label='Job Name'
                            icon='cube'
                            fluid
                            placeholder='Enter job name'
                            name='name'
                            value={name}
                            onChange={(e, { value }) => setName(value)}
                        />
                    </Grid.Column>

                    <Grid.Column>
                        <XTagTypeahead labeled onChange={(e, { value }) => setTags(value)} />
                    </Grid.Column>

                    <Grid.Column>
                        <XTargetTypeahead labeled onChange={(e, { value }) => setTargets(value)} />
                    </Grid.Column>
                </Grid>

                <Header inverted attached='top' size='large'>
                    <Icon name='code' />
                    {name ? name : 'Script'}
                </Header>
                <Form.Field
                    control={XScriptEditor}
                    content={content}
                    onChange={(e, { value }) => setContent(value)} />
                <Message
                    error
                    icon='warning'
                    header={'Failed to Queue Job'}
                    onDismiss={(e, data) => setError(null)}
                    content={err ? err.message : 'Unknown Error'}
                />
            </Modal.Content>
            <Modal.Actions>
                <Form.Button style={{ marginBottom: '10px' }} positive floated='right'>Queue</Form.Button>
            </Modal.Actions>
        </Modal >
    );
}

export default XJobQueueModal
