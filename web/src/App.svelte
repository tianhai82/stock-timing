<script>
  import { Button, NavigationDrawer, List, ListItem } from "smelte";
  import { onMount } from "svelte";
  import Router from "svelte-spa-router";
  import { location } from "svelte-spa-router";
  import Signin from "./components/Signin.svelte";
  import Stock from "./components/Stock.svelte";
  import Subscriptions from "./components/Subscriptions.svelte";
  import { loginUser, showSignIn, instruments, subscriptions } from "./store/store";
  import { retrieveSubscriptions, retrieveInstruments } from "./api/api";

  let showMenu = false;

  firebase.auth().onAuthStateChanged(function(user) {
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
      .catch(function(error) {
        console.error("sign out failed");
      });
  }

  const routes = {
    "/subscriptions": Subscriptions,
    "/:instrumentID?/:period?": Stock
  };

  let menu = [{ to: "#/", text: "Stock Analysis" }];
  $: {
    if (!!$loginUser) {
      menu = [
        { to: "#/", text: "Stock Analysis" },
        { to: "#/subscriptions", text: "Subscriptions" }
      ];
    } else {
      menu = [{ to: "#/", text: "Stock Analysis" }];
    }
  }

  function isSelected(url) {
    switch (url) {
      case "#/":
        if ($location.includes("subscriptions")) {
          return false;
        }
        return true;
      case "#/subscriptions":
        if ($location.includes("subscriptions")) {
          return true;
        }
        return false;
    }
    return false;
  }

  onMount(() => {
    retrieveInstruments().then(data => {
      $instruments = data.map(i => ({
        symbol: i.SymbolFull,
        value: i.InstrumentID,
        text: i.InstrumentDisplayName
      }));
    });
  });
</script>

<div class="h-auto">
  <NavigationDrawer
    bind:showDesktop={showMenu}
    bind:showMobile={showMenu}
    asideClasses="fixed top-0 h-full w-auto drawer overflow-hidden"
    breakpoint="sm">
    <h6 class="p-6 ml-1 pb-2 text-xs text-gray-900">Menu</h6>
    <List items={menu}>
      <span slot="item" let:item class="cursor-pointer">
        <a href={item.to}>
          <ListItem selected={isSelected(item.to)} {...item} navigation />
        </a>
      </span>
    </List>
  </NavigationDrawer>
  <header class="bg-teal-100 p-3 shadow">
    <div class="flex flex-row justify-between align-middle">
      <div class="flex flex-row">
        <Button
          class="m-0 mt-2 mr-2 p-0"
          color="black"
          icon="menu"
          text
          light
          flat
          on:click={() => (showMenu = !showMenu)} />
        <h5>Stock Timing</h5>
      </div>
      {#if $loginUser}
        <div class="flex flex-row align-middle">
          <span class="uppercase mr-2 mt-1">{$loginUser.displayName}</span>
          <Button
            on:click={signOut}
            dark
            icon="exit_to_app"
            text
            class="m-0 p-0" />
        </div>
      {:else}
        <Button on:click={() => ($showSignIn = true)}>Log In</Button>
      {/if}
    </div>
  </header>
  <Signin bind:showSignIn={$showSignIn} />
  <div class="container mx-auto items-center h-auto">
    <Router {routes} />
  </div>
</div>
