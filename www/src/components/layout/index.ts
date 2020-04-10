import './index.css';

export type XError = {
    title: string;
    msg: string;
}

export type XLoading = {
    title: string;
    msg: string;
}

export interface XViewProps {
    setLoading: (loading?: XLoading) => void;
    setError: (error?: XError) => void;
}

export { default as X404 } from './X404';
export { default as XBoundary } from './XBoundary';
export { default as XCardGroup } from './XCardGroup';
export { default as XLayout } from './XLayout';
export { default as XPrivateRoute } from './XPrivateRoute';
export { default as XSidebar } from './XSidebar';
export { default as XToolbar } from './XToolbar';
export { default as XUnimplemented } from './XUnimplemented';
export { default as XContent } from './XView';

