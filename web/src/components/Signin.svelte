<script>
import { Dialog, Button, Spinner } from 'svetamat';
import { auth } from '../fire.js';

export let showSignIn;
let promise;

function signInGoogle() {
  const provider = new firebase.auth.GoogleAuthProvider();
  promise = firebase
    .auth()
    .signInWithPopup(provider)
    .then(() => {
      showSignIn = false;
    })
    .catch(function (error) {
      var errorMessage = error.message;
      var email = error.email;
      alert(`Login failed. Account: ${email}. Error: ${errorMessage}`);
    });
}
</script>

<Dialog bind:visible={showSignIn}>
  <div class="p-6 bg-white w-40 rounded">
    <h6 class="mb-4 font-normal">Log in</h6>
    {#await promise}
      <Spinner/>
    {:then data}
      <Button on:click={signInGoogle} bgColor="bg-white">
        <div class="flex items-center">
          <img src="/images/google.svg" alt="Google" class="object-contain h-5"/>
          <span class="text-gray-900 ml-2 font-normal normal-case">Google</span>
        </div>
      </Button>
    {/await}
  </div>
</Dialog>
