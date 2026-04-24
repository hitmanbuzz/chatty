import { Component, Setter } from "solid-js";
import { MAX_USERNAME_LEN } from "../util/utils";

interface InputBox {
    text: string;
    setter: Setter<string>;
};

const Input: Component<InputBox> = (props) => {
    return (
        <input
            class="text-white border-none bg-gray-900 w-full max-w-[21rem] h-8 px-2"
            placeholder="Enter your username here"
            name="username-field"
            required
            maxLength={MAX_USERNAME_LEN}
            onInput={(e) => props.setter(e.currentTarget.value)}
        />
    );
}

export default Input
