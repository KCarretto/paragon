export type UserInfo = {
    userID?: number,
    authenticated: boolean,
    activated: boolean,
    admin: boolean,
}

export const getUserInfo = (): UserInfo => {
    let user : UserInfo = {
        authenticated: false,
        activated: false,
        admin: false,
    }

    let userJSON = localStorage.getItem("user");
    if (userJSON !== null) {
        localStorage.setItem("user", JSON.stringify(user));
        return user;
    }

    let userInfo : UserInfo = JSON.parse(userJSON);
    user.authenticated = userInfo.authenticated || false;
    user.activated = userInfo.activated || false;
    user.admin = userInfo.admin || false;
    if (userInfo.userID) {
        user.userID = userInfo.userID;
    }
    return user;
}
