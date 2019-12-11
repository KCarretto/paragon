import PropTypes from 'prop-types';
import React from 'react';
import { Link } from 'react-router-dom';
import { Icon, Menu, Sidebar } from 'semantic-ui-react';
import './index.css';

const XSidebar = (props) => (
    <Sidebar.Pushable className='XLayout'>
        <Sidebar
            as={Menu}
            icon='labeled'
            animation='push'
            direction='left'
            visible
            vertical
            inverted
            width='thin'
            className='XSidebar'
        >
            {props.routeMap ? props.routeMap.map((value, index) => {
                return <Menu.Item as={Link} to={value.link} key={index}>{value.icon}{value.title}</Menu.Item>
            }) : []}
            <Menu.Item href="https://github.com/kcarretto/paragon/issues/new"><Icon name='bug' />Bug</Menu.Item>
        </Sidebar>
        <Sidebar.Pusher className='XContent'>
            {props.children}
        </Sidebar.Pusher>
    </Sidebar.Pushable >
)
XSidebar.propTypes = {
    routeMap: PropTypes.arrayOf(
        PropTypes.shape({
            title: PropTypes.string.isRequired,
            link: PropTypes.string.isRequired,
            icon: PropTypes.element.isRequired,
            routes: PropTypes.arrayOf(PropTypes.element.isRequired),
        }).isRequired
    ).isRequired
}
export default XSidebar