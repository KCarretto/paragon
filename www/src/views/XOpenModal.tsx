import * as React from "react";
import { FunctionComponent } from "react";
import { Modal } from "semantic-ui-react";

type Openable = {
  open: boolean;
};

const XOpenModal: FunctionComponent<Openable> = props => (
  <Modal open={props.open}>{props.children}</Modal>
);

export default XOpenModal;
