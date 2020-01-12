import React from 'react';
import { Message } from 'semantic-ui-react';

export default ({ title, err }) => (
    <Message negative hidden={!err ? true : false}>
        <Message.Header>{title}</Message.Header>
        {err && err.message ? err.message : ''}
    </Message>
)

