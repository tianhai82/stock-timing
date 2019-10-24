const firebaseConfig = {
  apiKey: "AIzaSyBq6s254ST7NQAKr6oD78ojrJHjYdtqXjo",
  authDomain: "stock-timing.firebaseapp.com",
  databaseURL: "https://stock-timing.firebaseio.com",
  projectId: "stock-timing",
  storageBucket: "stock-timing.appspot.com",
  messagingSenderId: "249657550481",
  appId: "1:249657550481:web:45d6e4fc9276a220e54130"
};
const app = firebase.initializeApp(firebaseConfig);
const storage = app.storage();
const auth = app.auth();

export { app, storage, auth };
