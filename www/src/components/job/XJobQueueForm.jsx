import React from 'react';
import { Dropdown, Form } from 'semantic-ui-react';

export default class XJobQueueForm extends React.Component {
    state = {}

    render() {
        return (
            <Form>
                <Form.Dropdown text='Select a Job type...'>
                    <Dropdown.Menu>
                        <Dropdown.Item text='Command' />
                        <Dropdown.Item text='Teamserver' />
                        <Dropdown.Item text='Custom' />
                    </Dropdown.Menu>
                </Form.Dropdown>
                <Form.Input
                    icon='tags'
                    iconPosition='left'
                    label={{ content: 'Add Tag' }}
                    labelPosition='right'
                    placeholder='Add tags'
                />
                <Form.TextArea label='Task Content' placeholder='Script to execute...' />
                <Form.Button primary floated='right'>Queue</Form.Button>
            </Form>
        );
    }
}
