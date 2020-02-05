const firebaseConfig = {
  apiKey: "AIzaSyBq6s254ST7NQAKr6oD78ojrJHjYdtqXjo",
  authDomain: "stock-timing.firebaseapp.com",
  projectId: "stock-timing",
  appId: "1:249657550481:web:45d6e4fc9276a220e54130",
  measurementId: "G-J45TLZECXZ"
};
const app = firebase.initializeApp(firebaseConfig);
const auth = app.auth();
const analytics = app.analytics();
export { app, analytics, auth };
