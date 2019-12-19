import React from 'react';
import { Button, Icon, Modal } from 'semantic-ui-react';

const XJobCreateModal = () => (
    <Modal centered={false} trigger={
        <Button basic circular><Icon.Group size='big'>
            <Icon link fitted name='cube' color='green' />
            <Icon link fitted corner name='add' color='green' />
        </Icon.Group></Button>}>
        <Modal.Header>Create a Job</Modal.Header>
        <Modal.Content>
            <h1>Test</h1>
        </Modal.Content>
    </Modal>
)

XJobCreateModal.propTypes = {}

export default XJobCreateModal