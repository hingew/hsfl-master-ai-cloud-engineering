import { Elm } from './Main.elm'

const app = Elm.Main.init({
    node: document.getElementById('root'),
    flags: {
        //token: localStorage.getItem("token")
        token: "test-token"
    }
})

app.ports.storeToken.subscribe((token) => {
    localStorage.setItem("token", token)
    app.ports.gotToken.send(localStorage.getItem("token"))
})

