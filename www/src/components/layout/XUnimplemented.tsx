import React from 'react';
import { Grid, Message } from 'semantic-ui-react';

const XUnimplemented = (props) => (
    <Grid padded relaxed='very' centered>
        <Message
            icon='wrench'
            header='Under Construction'
            content={
                <span>This page has yet to be implemented, please try again later.
                    <a
                        href="https://github.com/KCarretto/paragon/pulls"
                        target="_blank"
                        rel="noopener noreferrer"
                    >
                        Want to help?
                        </a>
                </span>}
            size='massive'
            warning
        />
    </Grid>
)

export default XUnimplemented
