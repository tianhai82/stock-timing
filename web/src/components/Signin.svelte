<script>
  import { Dialog, Button, Image } from "smelte";
  import { tick } from "svelte";
  import { auth } from "../fire.js";

  export let showSignIn;
  function signInGoogle() {
    var provider = new firebase.auth.GoogleAuthProvider();
    firebase
      .auth()
      .signInWithPopup(provider)
      .then(() => {
        showSignIn = false;
      })
      .catch(function(error) {
        var errorMessage = error.message;
        var email = error.email;
        alert(`Login failed. Account: ${email}. Error: ${errorMessage}`);
      });
  }
</script>

<Dialog bind:value={showSignIn}>
  <h6 slot="title">Log in</h6>
  <Button on:click={signInGoogle} color="white" block>
    <div class="flex flex-row justify-between inline-block align-middle">
      <Image src="/images/google.svg" alt="Google" width={18} />
      <span class="text-gray-900 ml-2 font-normal normal-case">Google</span>
    </div>
  </Button>
</Dialog>
