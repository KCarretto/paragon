import moment from 'moment';
import React from 'react';
import { Link } from 'react-router-dom';
import { Divider, Feed, Header, Icon, List } from 'semantic-ui-react';
import { XTaskStatus } from '.';

const XTaskSummary = ({ tasks }) => (
    <Feed>
        <Header sub>Recent Tasks</Header>
        {tasks ? tasks.map((task, index) => (
            <Feed.Event key={index}>
                <Feed.Label>
                    <Icon fitted size='big' {...XTaskStatus.getStatus(task).icon} />
                </Feed.Label>
                <Feed.Content>
                    <Feed.Summary>
                        <Link to={'/jobs/' + task.job.id}><List.Header>{task.job.name}
                        </List.Header></Link>
                    </Feed.Summary>
                    <Feed.Extra text>
                        {XTaskStatus.getStatus(task).text}
                    </Feed.Extra>
                    <Feed.Meta>
                        Last Updated: {moment.unix(XTaskStatus.getTimestamp(task)).fromNow()}
                    </Feed.Meta>
                    <Divider />
                </Feed.Content>
            </Feed.Event>
        )) : <Header sub disabled>No recent tasks</Header>}
    </Feed>
);

export default XTaskSummary;
