import React from 'react';
import { Grid, Message } from 'semantic-ui-react';

const XUnimplemented = (props) => (
    <Grid padded relaxed='very' centered>
        <Message
            icon='wrench'
            header='Under Construction'
            content='This page has yet to be implemented, please try again later.'
            size='massive'
            warning
        />
    </Grid>
)

export default XUnimplemented
