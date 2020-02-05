import { writable } from 'svelte/store';

export const loginUser = writable(null);

export const showSignIn = writable(false);

export const instruments = writable([]);

export const subscriptions = writable([]);