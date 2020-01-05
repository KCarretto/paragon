import React from 'react';
import { Button, Modal } from 'semantic-ui-react';
import { XJobQueueForm } from '.';

const XJobQueueModal = ({ header }) => (
    <Modal centered={false} trigger={<Button positive circular icon='plus' />}>
        <Modal.Header>{header ? header : "Queue a Job"}</Modal.Header>
        <Modal.Content>
            <XJobQueueForm />
        </Modal.Content>
    </Modal >
)

XJobQueueModal.propTypes = {}

export default XJobQueueModal
