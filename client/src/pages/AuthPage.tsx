import { Component, createSignal } from "solid-js";
import { AuthButton } from "../components/Button";
import Input from "../components/InputBox";
import { AuthMethod } from "../auth/auth_method";
import { Authentication } from "../auth/authentication";

const Auth: Component = () => {
    const [username, SetUsername] = createSignal("");
    const authMethod = new Authentication();

    return (
        <div class="bg-black min-h-screen">
            <div class="flex flex-row items-center justify-center gap-1">
            </div>
            <div class="flex flex-col items-center justify-center min-h-screen gap-8 px-4">
                <div class="w-full flex justify-center">
                    <Input text={username()} setter={SetUsername} />
                </div>
                <div class="flex items-center justify-center gap-6">
                    <AuthButton text="Register" method={AuthMethod.REGISTER} auth={authMethod} username={username().trim()} />
                    <AuthButton text="Login" method={AuthMethod.LOGIN} auth={authMethod} username={username().trim()} />
                </div>
            </div>
        </div>
    );
}

export default Auth
