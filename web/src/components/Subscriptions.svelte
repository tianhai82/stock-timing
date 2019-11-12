<script>
  import { Button, ProgressCircular, Dialog } from "smelte";
  import { push } from "svelte-spa-router";
  import { retrieveSubscriptions, removeSubscription } from "../api/api";
  import { loginUser } from "../store/store";

  let promise;
  let items = [];
  let showDialog = false;
  let toRemove;
  function deleteSubscription(instrument) {
    return () => {
      toRemove = instrument;
      showDialog = true;
    };
  }
  function navigateTo(instrument) {
    return () => {
      push(`/${instrument.instrumentID}/${instrument.period}`);
    };
  }
  function confirmRemove() {
    $loginUser
      .getIdToken(true)
      .then(idToken =>
        removeSubscription({ idToken, instrumentID: toRemove.instrumentID })
      )
      .then(() => {
        items = items.filter(
          item => item.instrumentID !== toRemove.instrumentID
        );
        toRemove = undefined;
        showDialog = false;
      })
      .catch(err => alert(err));
  }
  $: {
    if (!!$loginUser) {
      promise = $loginUser
        .getIdToken(true)
        .then(idToken => retrieveSubscriptions(idToken))
        .then(data => {
          items = data;
        })
        .catch(err => alert(err));
    }
  }
</script>

<div class="h-auto">
  <h6 class="pt-4 py-2 px-2">Subscriptions</h6>
  {#await promise}
    <div class="px-4 pt-4">
      Loading...
      <ProgressCircular />
    </div>
  {:then}
    <div class="rounded h-full overflow-y-auto">
      <ul class="py-2 rounded">
        {#each items as item, i}
          <li
            class="hover:bg-gray-transDark relative overflow-hidden transition
            p-4 text-gray-700 flex items-center z-10 py-2">
            <div class="flex flex-row justify-between w-full">
              <div
                class="flex flex-col p-0 cursor-pointer"
                on:click={navigateTo(item)}>
                <div class="font-medium">{item.instrumentDisplayName}</div>
                <div class="text-gray-600 p-0 text-sm">
                  Period: {item.period}
                </div>
              </div>
              <Button
                icon="delete"
                text
                flat
                class="p-2"
                color="black"
                on:click={deleteSubscription(item)} />
            </div>
          </li>
        {/each}
      </ul>
    </div>
  {/await}
  <Dialog bind:value={showDialog}>
    <h5 slot="title">Remove {toRemove.instrumentDisplayName}?</h5>
    <div slot="actions">
      <Button text on:click={() => (showDialog = false)}>Cancel</Button>
      <Button text on:click={confirmRemove}>Confirm</Button>
    </div>
  </Dialog>
</div>
