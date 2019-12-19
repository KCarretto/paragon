import React from 'react';
<<<<<<< HEAD
import { Button, Modal } from 'semantic-ui-react';
import { XJobQueueForm } from '.';

const XJobQueueModal = ({ header }) => (
    <Modal centered={false} trigger={<Button primary compact basic icon='chevron right' />}>
        <Modal.Header>{header ? header : "Queue a Job"}</Modal.Header>
        <Modal.Content>
            <XJobQueueForm />
=======
import { Container, Icon, Modal } from 'semantic-ui-react';

const XJobQueueModal = (props) => (
    <Modal centered={false} trigger={<Container {...props} ><Icon name='plus square' /></Container>}>
        <Modal.Header>Create a Job</Modal.Header>
        <Modal.Content>
            <h1>Test</h1>
>>>>>>> Darn icons
        </Modal.Content>
    </Modal >
)

XJobQueueModal.propTypes = {}

export default XJobQueueModal
