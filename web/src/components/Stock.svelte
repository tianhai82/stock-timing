<script>
  import { Button, ProgressCircular, ProgressLinear, Slider } from "smelte";
  import Select from "./widgets/Select.svelte";
  import {
    retrieveCandles,
    retrieveSignals,
    addSubscription
  } from "../api/api";
  import CandleChart from "./CandleChart.svelte";
  import debounce from "../debounce";
  import {
    loginUser,
    showSignIn,
    instruments,
    subscriptions
  } from "../store/store";

  export let params = {};
  let stock;
  let candles;
  let signals;
  let freq = 48;
  let period;
  let showAnalyzing = false;
  let stockName;

  let loadChartPromise;

  $: {
    if (params.instrumentID && +params.instrumentID !== stock) {
      stock = +params.instrumentID;
      const stockFound = $instruments.find(ins => ins.value === stock);
      if (stockFound) {
        stockName = stockFound.text;
        signals = [];
        loadChartPromise = retrieveCandles(stock).then(data => {
          candles = data;
        });
        freqChanged();
      }
    }
  }
  $: {
    if (params.period && +params.period !== period) {
      period = +params.period;
      setFreq((period - 15) * 4);
      signals = [];
    }
  }

  function setFreq(x) {
    freq = x;
  }

  function subcribe() {
    if (!$loginUser) {
      alert("Please login to subscribe to alerts");
      return;
    }
    const stockFound = $instruments.find(ins => ins.value === stock);
    if (!stockFound) {
      alert("Please select a company/ETF to subscribe");
      return;
    }

    $loginUser.getIdToken(true).then(idToken => {
      addSubscription({
        idToken,
        instrument: {
          symbol: stockFound.symbol,
          instrumentID: stockFound.value,
          instrumentDisplayName: stockFound.text
        },
        period
      })
        .then(data => {
          $subscriptions = [
            ...$subscriptions,
            {
              symbol: stockFound.symbol,
              instrumentID: stockFound.value,
              instrumentDisplayName: stockFound.text,
              period
            }
          ];
          alert(
            `You are subscribed to trading signals for "${stockFound.text}"!`
          );
        })
        .catch(err => alert(err));
    });
  }

  const freqChanged = debounce(() => {
    if (stock) {
      showAnalyzing = true;
      const signalPromise = retrieveSignals(stock, period).then(data => {
        signals = data;
        showAnalyzing = false;
      });
    }
  }, 700);

  function stockChanged(e) {
    params = {};
    stock = e.detail;
    freq = 48;
    signals = [];
    freqChanged();
    loadChartPromise = retrieveCandles(stock).then(data => {
      candles = data;
    });
  }

  const filterStocks = (stock, inputValue) =>
    stock.text.toLowerCase().includes(inputValue) ||
    stock.symbol.toLowerCase().includes(inputValue);

  $: candleClass = !!candles && candles.length > 0 ? "px-4" : "hidden";
  $: {
    period = freq / 4 + 15;
    freqChanged();
  }
  $: periodLabel = `Period (${period})`;
</script>

<style>
  .hidden {
    display: none;
  }
</style>

<div class="px-4 pt-4 pb-2">
  <Select
    minChar={3}
    filter={filterStocks}
    bind:value={stock}
    on:change={stockChanged}
    selectedLabel={stockName}
    outlined
    autocomplete
    label="Enter Company Name"
    items={$instruments} />
</div>
{#await loadChartPromise}
  <div class="px-4">
    Analyzing...
    <ProgressCircular />
  </div>
{:then}
  <div class={candleClass}>
    {#if showAnalyzing}
      <div>
        <ProgressLinear />
      </div>
    {/if}
    <CandleChart {candles} {signals} />
    <div class="mb-5">
      <Slider label={periodLabel} bind:value={freq} step="4" />
    </div>
    <div class="my-2">
      {#if $loginUser}
        <Button block outline on:click={subcribe}>Subscribe</Button>
      {:else}
        <Button block outline on:click={() => ($showSignIn = true)}>
          Login to subscribe
        </Button>
      {/if}
    </div>
  </div>
{/await}
