import * as React from "react";
import { Button, Modal } from "semantic-ui-react";
import { XTagCreateForm } from ".";

type TagCreateModalParams = {
  openOnStart?: boolean;
};

const XTagCreateModal = ({ openOnStart }: TagCreateModalParams) => (
  <Modal
    open={openOnStart}
    centered={false}
    trigger={<Button positive circular icon="plus" />}
  >
    <Modal.Header>Create a Tag</Modal.Header>
    <Modal.Content>
      <XTagCreateForm />
    </Modal.Content>
  </Modal>
);

XTagCreateModal.propTypes = {};

export default XTagCreateModal;
