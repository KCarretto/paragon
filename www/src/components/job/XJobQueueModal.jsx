import React from 'react';
import { Container, Icon, Modal } from 'semantic-ui-react';

const XJobQueueModal = (props) => (
    <Modal centered={false} trigger={<Container {...props} ><Icon name='plus square' /></Container>}>
        <Modal.Header>Create a Job</Modal.Header>
        <Modal.Content>
            <h1>Test</h1>
        </Modal.Content>
    </Modal >
)

XJobQueueModal.propTypes = {}

export default XJobQueueModal
