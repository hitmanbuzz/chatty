import { Component } from "solid-js";
import { AuthInterface, AuthMethod } from "../auth/auth_method";

interface AuthButtonProps {
    text: string;
    username: string;
    method: AuthMethod;
    auth: AuthInterface;
};

function DoAuth(method: AuthMethod, auth: AuthInterface, username: string) {
    auth.SetUsername(username);
    
    switch (method) {
        case AuthMethod.REGISTER:
            auth.Register();
            break;
        case AuthMethod.LOGIN:
            auth.Login()
            break;
    }
}

export const AuthButton: Component<AuthButtonProps> = (props) => {
    return (
        <div>
            <button
                onclick={() => DoAuth(props.method, props.auth, props.username)}
                class="border rounded-sm px-2 py-1 bg-gray-800 text-white cursor-pointer active:scale-90 transition-transform border-transparent">
                {props.text}
            </button>
        </div>
    );
};
