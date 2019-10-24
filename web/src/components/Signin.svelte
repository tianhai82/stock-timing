<script>
  import { Dialog, Button, Image } from "smelte";
  import { tick } from "svelte";
  import { auth } from "../fire.js";

  export let showSignIn;
  function signIn() {
    var provider = new firebase.auth.GoogleAuthProvider();
    firebase
      .auth()
      .signInWithPopup(provider)
      .then(function(result) {
        // This gives you a Google Access Token. You can use it to access the Google API.
        var token = result.credential.accessToken;
        // The signed-in user info.
        var user = result.user;
        console.log(token, user);
        // ...
      })
      .catch(function(error) {
        // Handle Errors here.
        var errorCode = error.code;
        var errorMessage = error.message;
        // The email of the user's account used.
        var email = error.email;
        // The firebase.auth.AuthCredential type that was used.
        var credential = error.credential;
        console.log(errorCode);
        // ...
      });
  }
</script>

<Dialog bind:value={showSignIn}>
  <h5 slot="title">Sign in</h5>
  <Button on:click={signIn} color="pink" light block>
    <Image src="/images/google.svg" alt="Google" width={15}/>
    Sign in with Google
  </Button>
</Dialog>
