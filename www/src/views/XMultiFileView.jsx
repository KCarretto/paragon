import { useQuery } from '@apollo/react-hooks';
import gql from 'graphql-tag';
import React from 'react';
import { Card, Container, Loader, Menu } from 'semantic-ui-react';
import { XFileCard, XFileUploadModal } from '../components/file';
import { XErrorMessage } from '../components/messages';
export const MULTI_FILE_QUERY = gql`
{
	files {
        id
        name
        contentType
        size
        creationTime
        lastModifiedTime

        links {
        id
        alias
        clicks
        expirationTime
        }
    }
}`

const XMultiFileView = () => {
    const { called, loading, error, data } = useQuery(MULTI_FILE_QUERY, {
        pollInterval: 5000,
    });

    const showCards = () => {
        if (!data || !data.files || data.files.length < 1) {
            return (
                // TODO: Better styling
                <h1>No files found!</h1>
            );
        }
        return (<Card.Group centered itemsPerRow={4}>
            {data.files.map(file => (<XFileCard key={file.id} {...file} />))}
        </Card.Group>);
    };

    return (
        <Container style={{ padding: '10px' }}>
            <Menu secondary>
                <Menu.Item position='right'><XFileUploadModal /></Menu.Item>
            </Menu>
            <Container fluid style={{ padding: '20px' }}>
                <Loader disabled={!called || !loading} />
                <XErrorMessage title='Error Loading Files' err={error} />
                {showCards()}
            </Container>
        </Container>
    );
}

// XTargetCardGroup.propTypes = {
//     targets: PropTypes.arrayOf(PropTypes.shape({
//         id: PropTypes.number.isRequired,
//         name: PropTypes.string.isRequired,
//         tags: PropTypes.arrayOf(PropTypes.string),
//     })).isRequired,
// };

export default XMultiFileView