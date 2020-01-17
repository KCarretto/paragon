import React from 'react';
import { Icon } from 'semantic-ui-react';

export default ({ tags, defaultText }) => {
    if (!tags || tags.length < 1) {
        if (!defaultText) {
            return (<span />);
        }
        return (<span><Icon name='tags' />{defaultText}</span>);
    }
    return (
        <span><Icon name='tags' /> {tags.map(tag => tag.name).join(', ')}</span>
    );
}