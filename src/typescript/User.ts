import {Writable, writable} from "svelte/store";
const storageUserKey = "user"

export interface User {
    Name: string
    Admin: boolean
    Token: string
}

const user_string = window.localStorage.getItem(storageUserKey);
let user_value = null;
if (user_string !== null) {
    user_value = JSON.parse(user_string)
}

export const user: Writable<User | null> = writable(user_value);

user.subscribe(value => {
    window.localStorage.setItem(storageUserKey, JSON.stringify(value));
})