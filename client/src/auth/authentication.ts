import { MAX_USERNAME_LEN } from "../util/utils";
import { AuthInterface } from "./auth_method";

export class Authentication implements AuthInterface {
    private username: string = "";

    public async Register(): Promise<void> {
        if (this.username.length > MAX_USERNAME_LEN || this.username.length <= 0) {
            console.log("username is either empty or passed the max limit characters");
            return
        }
        
        const apiURL = import.meta.env.VITE_API_URL;
        const req = new Request(`${apiURL}/create-user`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                username: this.username
            }),
        });

        const response = await fetch(req);
        console.log("Registering....");
        console.log(await response.json());
    }

    public async Login(): Promise<void> {
        console.log(`Username: ${this.username}`);
        console.log("Logging....");
    }

    public SetUsername(username: string): void {
        this.username = username;
    }
} 
