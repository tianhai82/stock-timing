<script>
import Tailwindcss from './Tailwindcss.svelte';
import Button from './widgets/Button.svelte';
import NavigationDrawer from './widgets/NavigationDrawer.svelte';
import { onMount } from 'svelte';
import Router from 'svelte-spa-router';
import { location, push } from 'svelte-spa-router';
import Signin from './components/Signin.svelte';
import Stock from './components/Stock.svelte';
import Subscriptions from './components/Subscriptions.svelte';
import {
  loginUser,
  showSignIn,
  instruments,
  subscriptions,
} from './store/store';
import { retrieveSubscriptions, retrieveInstruments } from './api/api';

let showMenu = false;

firebase.auth()
  .onAuthStateChanged(function (user) {
    if (user) {
      $loginUser = user;
      fetchSubscriptions();
    } else {
      $loginUser = undefined;
    }
  });

function fetchSubscriptions() {
  $loginUser
    .getIdToken()
    .then(idToken => retrieveSubscriptions(idToken))
    .then(data => {
      $subscriptions = data;
    })
    .catch(err => alert(err));
}

function signOut() {
  firebase
    .auth()
    .signOut()
    .catch(function (error) {
      console.error('sign out failed');
    });
}

const routes = {
  '/subscriptions': Subscriptions,
  '/:instrumentID?/:period?/:buyFreq?/:sellFreq?': Stock,
};

let menu = [{
  to: '#/',
  text: 'Stock Analysis',
}];
$: {
  if (!!$loginUser) {
    menu = [
      {
        to: '#/',
        text: 'Stock Analysis',
      },
      {
        to: '#/subscriptions',
        text: 'Subscriptions',
      },
    ];
  } else {
    menu = [{
      to: '#/',
      text: 'Stock Analysis',
    }];
  }
}

function isSelected(url) {
  switch (url) {
    case '#/':
      if ($location.includes('subscriptions')) {
        return false;
      }
      return true;
    case '#/subscriptions':
      if ($location.includes('subscriptions')) {
        return true;
      }
      return false;
  }
  return false;
}

function goto(url) {
  push(url);
  showMenu = false;
}

onMount(() => {
  retrieveInstruments()
    .then(data => {
      $instruments = data.map(i => ({
        symbol: i.SymbolFull,
        value: i.InstrumentID,
        text: i.InstrumentDisplayName,
      }));
    });
});
</script>

<Tailwindcss/>
<div class="h-auto">
  <NavigationDrawer marginTop="mt-12"
                    bind:visible={showMenu}>
    <div class="w-56 bg-orange-100 h-full">
      <h3 class="font-medium px-6 pb-3 pt-4 tracking-wide text-gray-900">
        Menu
      </h3>
      <ul>
        {#each menu as item, i}
          <li on:click|preventDefault={goto(item.to)}
              class="px-4 py-3 hover:bg-gray-200 text-gray-800 text-sm tracking-wide cursor-pointer">
            {item.text}
          </li>
        {/each}
      </ul>
    </div>
  </NavigationDrawer>
  <header class="bg-blue-900 fixed left-0 right-0 top-0 h-12 mt-0 z-30 flex items-center
  justify-between">
    <div class="flex items-center">
      <i
        class="material-icons text-white ml-4 cursor-pointer ripple"
        on:click={() => (showMenu = !showMenu)}>
        menu
      </i>
      <img src="/images/TtT.svg" alt="logo" class="ml-3" style="height:28px;"/> <span
      class="ml-3 text-lg text-white font-medium">Time to Trade</span>
    </div>
    {#if $loginUser}
      <div class="flex items-center">
          <span class="uppercase mr-2 text-white">
            {$loginUser.displayName}
          </span>
        <i
          class="material-icons text-orange-300 cursor-pointer ripple mr-4"
          on:click={signOut}>
          exit_to_app
        </i>
      </div>
    {:else}
      <span class="mr-2">
        <Button
          on:click={() => ($showSignIn = true)}
          outlined
          outlineColor="border-orange-300"
          textColor="text-orange-300">
          Log In
        </Button>
      </span>
    {/if}
  </header>
  <Signin bind:showSignIn={$showSignIn}/>
  <div class="mt-12 container mx-auto items-center">
    <Router {routes}/>
  </div>
</div>
