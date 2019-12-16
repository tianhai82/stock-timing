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
  let buyFreq = 50;
  let sellFreq = 50;
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
  $: {
    if (params.buyFreq && !equalToBuyFreq(+params.buyFreq)) {
      setBuyFreq(+params.buyFreq);
    }
  }
  $: {
    if (params.sellFreq && !equalToSellFreq(+params.sellFreq)) {
      setSellFreq(+params.sellFreq);
    }
  }

  const equalToBuyFreq = f => f === buyFreq;
  const equalToSellFreq = f => f === sellFreq;

  function setBuyFreq(x) {
    buyFreq = x;
  }
  function setSellFreq(x) {
    sellFreq = x;
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
        period,
        buyLimit: (100 - buyFreq) / 200 + 0.25,
        sellLimit: (100 - sellFreq) / 200 + 0.25
      })
        .then(data => {
          $subscriptions = [
            ...$subscriptions.filter(s => s.symbol !== stockFound.symbol),
            {
              symbol: stockFound.symbol,
              instrumentID: stockFound.value,
              instrumentDisplayName: stockFound.text,
              period,
              buyLimit: (100 - buyFreq) / 200 + 0.25,
              sellLimit: (100 - sellFreq) / 200 + 0.25
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
      const signalPromise = retrieveSignals(
        stock,
        period,
        (100 - buyFreq) / 200 + 0.25,
        (100 - sellFreq) / 200 + 0.25
      ).then(data => {
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
  $: {
    if (buyFreq >= 0.0 || sellFreq >= 0.0) {
      freqChanged();
    }
  }
  $: periodLabel = `Period (${period})`;
  $: buyFreqLabel = `Buy Frequency (${buyFreq})`;
  $: sellFreqLabel = `Sell Frequency (${sellFreq})`;
</script>

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
    <div class="flex flex-wrap mb-5">
      <div class="w-full md:w-1/2 px-0 md:pr-2 mb-5 md:mb-0">
        <Slider label={buyFreqLabel} bind:value={buyFreq} step="1" />
      </div>
      <div class="w-full md:w-1/2 px-0 md:pl-2">
        <Slider label={sellFreqLabel} bind:value={sellFreq} step="1" />
      </div>
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
