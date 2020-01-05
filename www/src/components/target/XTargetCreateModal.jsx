import React from 'react';
import { Button, Modal } from 'semantic-ui-react';
import { XTargetCreateForm } from '.';

const XTargetCreateModal = () => (
    <Modal centered={false} trigger={<Button positive circular icon='plus' />}>
        <Modal.Header>Create a Target</Modal.Header>
        <Modal.Content>
            <XTargetCreateForm />
        </Modal.Content>
    </Modal >
)

XTargetCreateModal.propTypes = {}

export default XTargetCreateModal
