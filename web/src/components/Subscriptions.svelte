<script>
import { Button, Dialog, Spinner } from 'svelte-tailwind-material';
import { push } from 'svelte-spa-router';
import { removeSubscription } from '../api/api';
import { loginUser, subscriptions } from '../store/store';

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
    push(
      `/${instrument.instrumentID}/${instrument.period}/${getFreq(
        instrument.buyLimit,
      )}/${getFreq(instrument.sellLimit)}`,
    );
  };
}

function confirmRemove() {
  $loginUser
    .getIdToken(true)
    .then(idToken =>
      removeSubscription({
        idToken,
        instrumentID: toRemove.instrumentID,
      }),
    )
    .then(() => {
      $subscriptions = $subscriptions.filter(
        item => item.instrumentID !== toRemove.instrumentID,
      );
      toRemove = undefined;
      showDialog = false;
    })
    .catch(err => alert(err));
}

function getFreq(limit) {
  if (limit === 0.0) {
    return 50;
  }
  return Math.round(100 - (limit - 0.25) * 200);
}
</script>

<div class="h-auto overflow-hidden z-0">
  <h6 class="pt-4 pb-1 px-2 font-medium">Subscriptions</h6>
  <div class="rounded h-full">
    <ul class="py-2 rounded">
      {#each $subscriptions as item, i}
        <li
          class="hover:bg-gray-300 relative overflow-hidden transition z-0
            p-4 text-gray-800 flex items-center z-10 py-2">
          <div class="flex flex-row justify-between w-full">
            <div
              class="flex flex-col p-0 cursor-pointer"
              on:click={navigateTo(item)}>
              <div class="font-medium">{item.instrumentDisplayName}</div>
              <div class="text-gray-700 p-0 text-sm">
                Period: {item.period} | Buy Frequency: {getFreq(item.buyLimit)}
                | Sell Frequency: {getFreq(item.sellLimit)}
              </div>
            </div>
            <i
              class="material-icons text-black p-2 cursor-pointer ripple"
              on:click={deleteSubscription(item)}>
              delete
            </i>
          </div>
        </li>
      {/each}
    </ul>
  </div>
  <Dialog bind:visible={showDialog}>
    <div class="px-4 py-5 bg-white rounded shadow-lg">
      <h3 class="text-xl">Remove {toRemove.instrumentDisplayName}?</h3>
      <div class="flex justify-end mt-6">
        <Button text textColor="text-orange-600" on:click={() => (showDialog = false)}>
          Cancel
        </Button>
        <Button outlined textColor="text-orange-600" outlineColor="border-orange-600"
                on:click={confirmRemove}>Confirm
        </Button>
      </div>
    </div>
  </Dialog>
</div>
