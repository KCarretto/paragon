
const XCredentialStatus = {
    WORKING: { text: 'The Credential has never failed.', icon: { name: 'check circle', color: 'green', className: 'XCircleIcon', bordered: false, circular: true } },

    UNSURE: { text: 'The credential has failed atleast once.', icon: { name: 'times circle', color: 'yellow', className: 'XCircleIcon', bordered: false, circular: true } },

    FAILING: {
        text: 'The credential has failed atleast eight times.', icon: { name: 'times circle', color: 'red', className: 'XCircleIcon', bordered: false, circular: true }
    },
}

XCredentialStatus.getStatus = ({ principal, secret, fails }) => {
    if (fails > 8) {
        return XCredentialStatus.FAILING;
    } else if (fails > 0) {
        return XCredentialStatus.UNSURE;
    }
    return XCredentialStatus.WORKING;
}

export default XCredentialStatus;