import PropTypes from 'prop-types';
import React from 'react';
import { BrowserRouter as Router, Switch } from 'react-router-dom';
import XSidebar from './XSidebar';


const XLayout = (props) => (
    <Router>
        <Switch>
            <XSidebar routeMap={props.routeMap}>
                {props.routeMap ? props.routeMap.map((value, index) => {
                    return <div route_group={value.title} key={index}>
                        {value.routes}
                    </div>
                }) : []}
                {props.children}
            </XSidebar>
        </Switch>
    </Router>

)

XLayout.propTypes = {
    routeMap: PropTypes.arrayOf(
        PropTypes.shape({
            title: PropTypes.string.isRequired,
            link: PropTypes.string.isRequired,
            icon: PropTypes.element.isRequired,
            routes: PropTypes.arrayOf(PropTypes.element.isRequired),
        }).isRequired
    ).isRequired
}
export default XLayout