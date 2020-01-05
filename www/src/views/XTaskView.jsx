import { useQuery } from '@apollo/react-hooks';
import gql from 'graphql-tag';
import React from 'react';
import { useParams } from 'react-router-dom';
import { Card, Container, Icon, Label, Loader } from 'semantic-ui-react';
import { XTaskStatus } from '../components/task';


const TASK_QUERY = gql`
    query Task($id: ID!) {
    task(id: $id) {
        id
        queueTime
        claimTime
        execStartTime
        execStopTime
        content
        output
        error
        sessionID
        job {
            id
            name
            tags {
                id
                name
            }
        }
    }
  }
`;

const XTaskView = () => {
    let { id } = useParams();

    const { loading, error, data } = useQuery(TASK_QUERY, {
        variables: { id },
    });

    if (loading) return (<Loader active />);
    if (error) return (`${error}`);

    let status = XTaskStatus.getStatus(data.task).icon;

    return (
        <Container fluid style={{ padding: '20px' }}>
            <Card fluid centered>
                <Card.Content>
                    <Card.Header href={'/jobs/' + data.task.job.id}>
                        <Label size='mini' icon={<Icon fitted size='big' {...status} />} attached='top right' />
                        {data.task.job.name ? data.task.job.name : 'Untitled Job'}
                    </Card.Header>
                    {
                        data.task.job.tags && data.task.job.tags.length != 0 ?
                            <Card.Meta>
                                <Icon name='tags' /> {data.task.job.tags.map(tag => tag.name).join(', ')}
                            </Card.Meta>
                            :
                            <Card.Meta>
                                <Icon name='tags' /> None
                            </Card.Meta>
                    }
                    <Card.Description>
                        <div className="ui visible message">
                            <div className="header">
                                Content
                            </div>
                            <pre>{data.task.content}</pre>
                        </div>
                        {data.task.output ?
                            <div className="ui positive message">
                                <div className="header">
                                    Output
                                </div>
                                <pre>{data.task.output}</pre>
                            </div> :
                            <div></div>
                        }
                        {data.task.error ?
                            <div className="ui negative message">
                                <div className="header">
                                    Error
                                </div>
                                <pre>{data.task.error}</pre>
                            </div>
                            :
                            <div></div>
                        }
                    </Card.Description>
                </Card.Content>
                {data.task.sessionID ?
                    <Card.Content extra>
                        <a>
                            <i aria-hidden="true" className="user icon"></i>
                            {data.task.sessionID}
                        </a>
                    </Card.Content>
                    :
                    <div></div>
                }

            </Card>
        </Container>
    );
}

export default XTaskView;
