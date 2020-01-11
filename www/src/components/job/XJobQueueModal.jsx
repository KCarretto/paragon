import { useMutation } from '@apollo/react-hooks';
import gql from 'graphql-tag';
import React, { useState } from 'react';
import { ToastContainer } from 'react-toastr';
import { Button, Form, Loader, Modal } from 'semantic-ui-react';
import { MULTI_JOB_QUERY, MULTI_TARGET_QUERY } from '../../views';
import { XTagTypeahead, XTargetTypeahead } from '../form';

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
        });
    }

    const getModalContent = () => {
        if (called && loading) {
            return (
                <Modal.Content>
                    <Loader />
                </Modal.Content>
            );
        }
        return (
            <Modal.Content>
                <ToastContainer
                    ref={ref => container = ref}
                    className="toast-top-right"
                />
                <Form.Input
                    label={{ content: 'Name' }}
                    placeholder='Enter job name'
                    name='name'
                    value={params.name}
                    onChange={handleChange}
                />
                <XTargetTypeahead
                    onChange={handleChange}
                    selectedValues={params.targets}
                />
                <XTagTypeahead
                    onChange={handleChange}
                    selectedValues={params.tags}
                />
                <Form.TextArea
                    label={{ content: 'Enter script' }}
                    placeholder='Enter script content'
                    name='content'
                    rows={15}
                    value={params.content}
                    onChange={handleChange}
                />

            </Modal.Content>
        );
    }


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
            <Modal.Header>{header ? header : "Queue a Job"}</Modal.Header>
            {getModalContent()}
            <Modal.Actions>
                <Form.Button positive floated='right'>Queue</Form.Button>
            </Modal.Actions>
        </Modal >
    );
}

export default XJobQueueModal
