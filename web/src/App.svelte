<script>
  import { Button } from "smelte";
  import Select from "./components/widgets/Select.svelte";
  import Signin from "./components/Signin.svelte";
  import {
    retrieveInstruments,
    retrieveCandles,
    retrieveSignals
  } from "./api/api";
  import CandleChart from "./components/CandleChart.svelte";

  export let name;
  let showSignIn = false;
  let stock;
  let instruments;
  let candles;
  let signals;

  retrieveInstruments().then(data => {
    instruments = data.map(i => ({
      symbol: i.SymbolFull,
      value: i.InstrumentID,
      text: i.InstrumentDisplayName
    }));
  });
  function stockChanged(e) {
    stock = e.detail;
    retrieveCandles(stock).then(data => {
      candles = data;
    });
    retrieveSignals(stock).then(data => {
      signals = data;
    });
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
  <header class="bg-teal-200 p-3 shadow">
    <div class="flex flex-row align-middle">
      <img
        src="/images/time-money.png"
        alt="logo"
        class="object-contain h-8 mx-4 mt-1" />
      <h5>Stock Timing</h5>
      <Button on:click={() => (showSignIn = true)}>Sign In</Button>
    </div>
  </header>
  <Signin bind:showSignIn />
  <div class="px-4 pt-4 pb-2">
    <Select
      filter={filterStocks}
      bind:value={stock}
      on:change={stockChanged}
      outlined
      autocomplete
      label="Enter Company Name"
      items={instruments} />
  </div>
  <div class={candleClass}>
    <CandleChart {candles} {signals} />
  </div>

</div>
