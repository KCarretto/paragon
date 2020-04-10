import * as React from "react";
import { Redirect, Route, RouteComponentProps, RouteProps } from "react-router-dom";
import { Container, Icon, Message } from "semantic-ui-react";
import { XError } from ".";

const XView: React.FC<RouteProps & { authorized: boolean, padded?: boolean, view: React.ElementType }> = ({ authorized, padded, view, exact, path, children, ...props }) => {
    const [error, setError] = React.useState<XError>(null);

    const unauthorizedView: React.FC<RouteComponentProps> = ({ location }) => (
        <Redirect to={{
            pathname: "/login",
            state: { from: location }
        }} />
    );

    const errorView = () => (
        <Container style={{ marginTop: "25px" }}>
            <Message negative>
                <Icon name="times circle" />
                <Message.Header>{error.title || "Fatal Error"}</Message.Header>
                {error.msg || "¯\_(ツ)_/¯ This is probably a bug"}
            </Message>
        </Container>
    );

    const componentView = () => (
        <div className={padded ? "XView XPaddedView" : "XView"}>
            {React.createElement(
                view,
                {
                    setError: (err?: XError) => setError(err),
                    ...props,
                },
            )}
        </div>
    );

    if (!authorized) {
        return <Route exact={exact} path={path} render={unauthorizedView} />;
    }

    if (error) {
        return <Route exact={exact} path={path} render={errorView} />
    }

    // if (loading) {
    //     console.log("LOADING: ", loading)
    //     return <Route exact={exact} path={path} render={loadingView} />
    // }

    return <Route exact={exact} path={path} render={componentView} />
}

export default XView;