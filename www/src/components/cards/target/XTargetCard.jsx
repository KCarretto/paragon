import moment from 'moment';
import PropTypes from 'prop-types';
import React from 'react';
import { Link } from 'react-router-dom';
import { Card, Feed, Header, Icon, Label, List, Segment } from 'semantic-ui-react';
import { XTaskStatus } from '../task';

const XTargetCard = ({ id, name, primaryIP, lastSeen, tags, tasks }) => (
    <Card>
        <Card.Content>
            <Card.Header>{name} </Card.Header>
            <Label corner='right' size='large' icon='check circle' color='green' />
            <Card.Meta>{moment.unix(lastSeen).fromNow()}</Card.Meta>
            <Segment >
                <Feed>
                    <Header sub>Recent Tasks</Header>
                    {tasks ? tasks.map((task, index) => (
                        <Feed.Event key={index}>
                            <Feed.Label>
                                <Icon fitted size='big' {...XTaskStatus.getStatus(task).icon} />
                            </Feed.Label>
                            <Feed.Content>
                                <Feed.Summary>
                                    <Link to={'/jobs/' + task.job.id}>
                                        <List.Header>{task.job.name}
                                            {/* <Progress color='red' size='small' percent={50} active>In Progress</Progress> */}
                                        </List.Header>
                                    </Link>
                                </Feed.Summary>
                                <Feed.Extra text>
                                    {XTaskStatus.getStatus(task).text}
                                </Feed.Extra>
                                <Feed.Meta>
                                    Last Updated: {moment.unix(XTaskStatus.getTimestamp(task)).fromNow()}
                                </Feed.Meta>
                            </Feed.Content>
                        </Feed.Event>
                    )) : <Header sub disabled>No recent tasks</Header>}
                </Feed>
            </Segment>
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