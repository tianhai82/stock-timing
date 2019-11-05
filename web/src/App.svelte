<script>
  import { Button, ProgressCircular } from "smelte";
  import Select from "./components/widgets/Select.svelte";
  import Signin from "./components/Signin.svelte";
  import {
    retrieveInstruments,
    retrieveCandles,
    retrieveSignals,
    addSubscription
  } from "./api/api";
  import CandleChart from "./components/CandleChart.svelte";

  let loginUser;

  let showSignIn = false;
  let stock;
  let instruments;
  let candles;
  let signals;

  firebase.auth().onAuthStateChanged(function(user, x) {
    if (user) {
      // var displayName = user.displayName;
      // var email = user.email;
      // var emailVerified = user.emailVerified;
      // var photoURL = user.photoURL;
      // var isAnonymous = user.isAnonymous;
      // var uid = user.uid;
      // var providerData = user.providerData;
      loginUser = user;
    } else {
      loginUser = undefined;
    }
  });
  function signOut() {
    firebase
      .auth()
      .signOut()
      .catch(function(error) {
        console.error("sign out failed");
      });
  }

  function subcribe() {
    if (!loginUser) {
      alert("Please login to subscribe to alerts");
      return;
    }
    const stockFound = instruments.find(ins => ins.value === stock);
    if (!stockFound) {
      alert("Please select a company/ETF to subscribe");
      return;
    }

    loginUser.getIdToken(true).then(idToken => {
      addSubscription({
        idToken,
        instrument: {
          symbol: stockFound.symbol,
          instrumentID: stockFound.value,
          instrumentDisplayName: stockFound.text
        }
      })
        .then(data => alert(`You are subscribed to trading signals for "${stockFound.text}"!`))
        .catch(err => alert(err));
    });
  }

  retrieveInstruments().then(data => {
    instruments = data.map(i => ({
      symbol: i.SymbolFull,
      value: i.InstrumentID,
      text: i.InstrumentDisplayName
    }));
  });

  let loadChartPromise;

  function stockChanged(e) {
    stock = e.detail;
    const candlePromise = retrieveCandles(stock).then(data => {
      candles = data;
    });
    const signalPromise = retrieveSignals(stock).then(data => {
      signals = data;
    });
    loadChartPromise = Promise.all([candlePromise, signalPromise]);
  }
  const filterStocks = (stock, inputValue) =>
    stock.text.toLowerCase().includes(inputValue) ||
    stock.symbol.toLowerCase().includes(inputValue);

  $: candleClass = !!candles && candles.length > 0 ? "px-4" : "hidden";
</script>

<style>
  .hidden {
    display: none;
  }
</style>

<div class="container mx-auto h-full items-center">
  <header class="bg-teal-100 p-3 shadow">
    <div class="flex flex-row justify-between align-middle">
      <div class="flex flex-row">
        <img
          src="/images/time-money.png"
          alt="logo"
          class="object-contain h-8 mx-4 mt-1" />
        <h5>Stock Timing</h5>
      </div>
      {#if loginUser}
        <div class="flex flex-row align-middle">
          <span class="uppercase mr-2 mt-1">{loginUser.displayName}</span>
          <Button
            on:click={signOut}
            dark
            icon="exit_to_app"
            text
            class="m-0 p-0" />
        </div>
      {:else}
        <Button on:click={() => (showSignIn = true)}>Log In</Button>
      {/if}
    </div>
  </header>
  <Signin bind:showSignIn />
  <div class="px-4 pt-4 pb-2">
    <Select
      minChar={3}
      filter={filterStocks}
      bind:value={stock}
      on:change={stockChanged}
      outlined
      autocomplete
      label="Enter Company Name"
      items={instruments} />
  </div>
  {#await loadChartPromise}
    <div class="px-4">
      Analyzing...
      <ProgressCircular />
    </div>
  {:then}
    <div class={candleClass}>
      <CandleChart {candles} {signals} />
      <div class="mt-2">
        {#if loginUser}
          <Button block outline on:click={subcribe}>Subscribe</Button>
        {:else}
          <Button block outline on:click={() => (showSignIn = true)}>
            Login to subscribe
          </Button>
        {/if}
      </div>
    </div>
  {/await}
</div>
