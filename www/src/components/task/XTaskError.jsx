import React from 'react';
import { Header, Icon, Segment } from 'semantic-ui-react';

export default ({ error }) => {
    if (!error) {
        return (<span />);
    }

    return (
        <div>
            <Header inverted size='large' attached='top' color='red' style={{ marginTop: '10px' }}>
                <Icon name='warning circle' />
                <Header.Content>Error</Header.Content>
            </Header>
            <Segment raised attached>
                <pre>{error}</pre>
            </Segment>
        </div>
    );
}