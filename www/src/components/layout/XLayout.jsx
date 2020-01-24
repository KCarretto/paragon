import PropTypes from 'prop-types';
import React from 'react';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import XLogin from '../../views/XLogin';
import XSidebar from './XSidebar';


const XLayout = (props) => (
    <Router>
        <Switch>
            <Route
                path="/login"
                component={XLogin}
            />
            <Route
                path="/"
            >
                <XSidebar routeMap={props.routeMap}>
                    {props.children}
                </XSidebar>
            </Route>
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