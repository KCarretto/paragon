import React from 'react';
import { Button, Modal } from 'semantic-ui-react';

const XJobCreateModal = () => (
    <Modal centered={false} trigger={<Button positive circular icon='plus' />}>
        <Modal.Header>Create a Job</Modal.Header>
        <Modal.Content>
            <h1>Test</h1>
        </Modal.Content>
    </Modal>
)

XJobCreateModal.propTypes = {}

export default XJobCreateModal