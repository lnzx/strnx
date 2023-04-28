const TOKEN_KEY = "token";

function internalToken() {
    const token = sessionStorage.getItem(TOKEN_KEY) || localStorage.getItem(TOKEN_KEY)
    return token == null ? {} : JSON.parse(token);
}

const useSession = {
    setToken: (token, remember) => {
        const data = JSON.stringify(token)
        sessionStorage.setItem(TOKEN_KEY, data)
        if(remember){
            localStorage.setItem(TOKEN_KEY, data)
        }
    },

    getToken: () => { 
        return internalToken();
    },

    clearToken: () => {
        localStorage.removeItem(TOKEN_KEY)
        sessionStorage.removeItem(TOKEN_KEY)
    },

    isLogin: () => {
        const t = internalToken()
        return t.token && t.token != ''
    }
}

export default useSession;