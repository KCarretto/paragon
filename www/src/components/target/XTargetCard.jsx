import moment from 'moment';
import PropTypes from 'prop-types';
import React from 'react';
import { Card, Icon, Label } from 'semantic-ui-react';
import { XTaskSummary } from '../task';

const XTargetCard = ({ id, name, primaryIP, lastSeen, tags, tasks }) => (
    <Card fluid >
        <Card.Content>
            <Card.Header href={'/targets/' + id}>{name} </Card.Header>
            {
                (!lastSeen || moment(lastSeen).isBefore(moment().subtract(5, 'minutes'))) ?
                    <Label corner='right' size='large' icon='times circle' color='red' />
                    : <Label corner='right' size='large' icon='check circle' color='green' />
            }
            <Card.Meta>{lastSeen ? moment(lastSeen).fromNow() : 'Never'}</Card.Meta>
            <XTaskSummary tasks={tasks} />
        </Card.Content>
        <Card.Content extra>
            <Icon name='tags' /> {tags ? tags.map(tag => tag.name).join(', ') : 'None'}
        </Card.Content>
    </Card>
)

XTargetCard.propTypes = {
    id: PropTypes.number.isRequired,
    name: PropTypes.string.isRequired,
    primaryIP: PropTypes.string,
    lastSeen: PropTypes.number,
    tags: PropTypes.arrayOf(PropTypes.shape({
        id: PropTypes.number.isRequired,
        name: PropTypes.string.isRequired
    })),
    tasks: PropTypes.arrayOf(PropTypes.shape({
        id: PropTypes.number.isRequired,

        queueTime: PropTypes.number.isRequired,
        claimTime: PropTypes.number,
        execStartTime: PropTypes.number,
        execStopTime: PropTypes.number,
        error: PropTypes.string,

        job: PropTypes.shape({
            id: PropTypes.number.isRequired,
            name: PropTypes.string.isRequired,
        }).isRequired
    }))
}

export default XTargetCard