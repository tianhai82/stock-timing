<script>
  import { Select, Button } from "smelte";
  import Signin from "./components/Signin.svelte";
  import { retrieveInstruments, retrieveCandles } from "./api/api";
  import CandleChart from "./components/CandleChart.svelte";

  export let name;
  let showSignIn = false;
  let stock;
  let instruments;
  let candles;

  retrieveInstruments().then(data => {
    instruments = data.map(i => ({
      value: i.InstrumentID,
      text: i.InstrumentDisplayName
    }));
  });
  function stockChanged() {
    retrieveCandles(stock).then(data => {
      candles = data;
    });
  }
</script>

<div class="container mx-auto h-full items-center">
  <header class="bg-teal-200 p-3 shadow">
    <div class="flex flex-row align-middle">
      <img
        src="/images/time-money.png"
        alt="logo"
        class="object-contain h-10 mx-4 mt-1" />
      <h4>Stock Timing</h4>
      <Button on:click={() => (showSignIn = true)}>Sign In</Button>
    </div>
  </header>
  <Signin bind:showSignIn />
  <div class="p-4">
    <Select
      bind:value={stock}
      on:change={stockChanged}
      outlined
      autocomplete
      label="Enter Company Name"
      items={instruments} />
    <CandleChart {candles} />
  </div>

</div>
