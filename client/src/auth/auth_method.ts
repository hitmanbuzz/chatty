export enum AuthMethod { REGISTER, LOGIN };

export interface AuthInterface {
   Register(): Promise<void>;
   Login(): Promise<void>; 
   SetUsername(username: string): void;
};
