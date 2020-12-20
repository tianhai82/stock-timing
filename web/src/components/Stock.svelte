<script>
  import {
    Button,
    Slider,
    Autocomplete,
    Spinner,
    Progress,
    Input,
  } from "svetamat";
  import {
    retrieveCandles,
    retrieveSignals,
    addSubscription,
  } from "../api/api";
  import CandleChart from "./CandleChart.svelte";
  import debounce from "../debounce";
  import {
    loginUser,
    showSignIn,
    instruments,
    subscriptions,
  } from "../store/store";

  export let params = {};
  let stock = { value: "" };
  let candles;
  let signals;
  let freq = 27;
  let buyFreq = 50;
  let sellFreq = 50;
  let showAnalyzing = false;
  let stockName;

  let loadChartPromise;

  $: {
    if (params.instrumentID) {
      const stockFound = $instruments.find(
        (ins) => ins.value === +params.instrumentID
      );
      if (stockFound) {
        stock = stockFound;
        stockName = stockFound.text;
        signals = [];
        loadChartPromise = retrieveCandles(stock.value).then((data) => {
          candles = data;
        });
        freqChanged();
      }
    }
  }
  $: {
    if (params.period && !equalToPeriod(+params.period)) {
      setFreq(+params.period);
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

  const equalToPeriod = (f) => f === freq;
  const equalToBuyFreq = (f) => f === buyFreq;
  const equalToSellFreq = (f) => f === sellFreq;

  function setFreq(x) {
    freq = x;
  }

  function setBuyFreq(x) {
    buyFreq = x;
  }

  function setSellFreq(x) {
    sellFreq = x;
  }

  function subcribe() {
    if (!$loginUser) {
      alert("Please login to subscribe to alerts");
      return;
    }
    const stockFound = $instruments.find((ins) => ins.value === stock.value);
    if (!stockFound) {
      alert("Please select a company/ETF to subscribe");
      return;
    }

    $loginUser.getIdToken(true).then((idToken) => {
      addSubscription({
        idToken,
        instrument: {
          symbol: stockFound.symbol,
          instrumentID: stockFound.value,
          instrumentDisplayName: stockFound.text,
        },
        period: freq,
        buyLimit: (100 - buyFreq) / 200 + 0.25,
        sellLimit: (100 - sellFreq) / 200 + 0.25,
      })
        .then((data) => {
          $subscriptions = [
            ...$subscriptions.filter((s) => s.symbol !== stockFound.symbol),
            {
              symbol: stockFound.symbol,
              instrumentID: stockFound.value,
              instrumentDisplayName: stockFound.text,
              period: freq,
              buyLimit: (100 - buyFreq) / 200 + 0.25,
              sellLimit: (100 - sellFreq) / 200 + 0.25,
            },
          ];
          alert(
            `You are subscribed to trading signals for "${stockFound.text}"!`
          );
        })
        .catch((err) => alert(err));
    });
  }

  const freqChanged = debounce(() => {
    if (stock && stock.value) {
      showAnalyzing = true;
      const signalPromise = retrieveSignals(
        stock.value,
        freq,
        (100 - buyFreq) / 200 + 0.25,
        (100 - sellFreq) / 200 + 0.25
      ).then((data) => {
        signals = data;
        showAnalyzing = false;
      });
    }
  }, 700);

  function stockChanged(e) {
    params = {};
    stock = e.detail;
    freq = 27;
    signals = [];
    freqChanged();
    loadChartPromise = retrieveCandles(stock.value).then((data) => {
      candles = data;
    });
  }

  $: candleClass = !!candles && candles.length > 0 ? "px-4" : "hidden";

  $: {
    if (buyFreq >= 0.0 || sellFreq >= 0.0 || (freq >= 15 && freq <= 40)) {
      freqChanged();
    }
  }
  $: periodLabel = `Period (${freq})`;
  $: buyFreqLabel = `Buy Frequency (${buyFreq})`;
  $: sellFreqLabel = `Sell Frequency (${sellFreq})`;

  let increase = 0.0;
  $: if (candles && candles.length > 2) {
    increase = (
      (candles[candles.length - 1].Close / candles[0].Open - 1) *
      100
    ).toFixed(2);
  } else {
    increase = 0.0;
  }

  let buyIncrease = 0.0;
  $: if (signals && signals.length > 0 && candles && candles.length > 1) {
    let buyAmount = 0.0;
    let currentAmount = 0.0;
    signals.forEach((s) => {
      if (s.Signal === 1) {
        buyAmount += 2000;
        const endPrice = candles[candles.length - 1].Close;
        const changePercent = endPrice / s.Price;
        currentAmount += changePercent * 2000;
      }
    });
    buyIncrease = ((currentAmount / buyAmount - 1) * 100).toFixed(2);
  } else {
    buyIncrease = 0.0;
  }

  let buySellIncrease = 0.0;
  $: if (signals && signals.length > 0 && candles && candles.length > 1) {
    let buyAmount = 0.0;
    let sellAmount = 0.0;
    let currentAmount = 0.0;
    signals.forEach((s) => {
      if (s.Signal === 1) {
        buyAmount += 2000;
        const sellPrice = candles[candles.length - 1].Close;
        const changePercent = sellPrice / s.Price;
        currentAmount += changePercent * 2000;
      } else if (s.Signal === 2) {
        sellAmount += 2000;
        const buyPrice = candles[candles.length - 1].Close;
        const changePercent = s.Price / buyPrice;
        currentAmount += changePercent * 2000;
      }
    });
    buySellIncrease = (
      (currentAmount / (buyAmount + sellAmount) - 1) *
      100
    ).toFixed(2);
  } else {
    buySellIncrease = 0.0;
  }

  let weightedPeriodAvg;
  $: if (candles && candles.length > 2) {
    let count = candles.length * 4;
    let total = 0;
    let p = 0;
    candles.forEach((candle, i) => {
      p += candle.Close * ((i + 1) / count);
      total += (i + 1) / count;
    });
    weightedPeriodAvg = (p / total).toFixed(2);
  } else {
    weightedPeriodAvg = 0.0;
  }

  let winningRate;
  $: if (signals && signals.length > 0 && candles && candles.length > 1) {
    let wins = 0;
    let losses = 0;
    signals.forEach((s) => {
      if (s.Signal === 1) {
        if (s.Price < weightedPeriodAvg) {
          wins += 1;
        }
        if (s.Price > weightedPeriodAvg) {
          losses += 1;
        }
      }
      // else if (s.Signal === 2) {
      //   if (s.Price < weightedPeriodAvg) {
      //     losses += 1;
      //   }
      //   if (s.Price > weightedPeriodAvg) {
      //     wins += 1;
      //   }
      // }
    });
    winningRate = ((wins / (wins + losses)) * 100).toFixed(2);
  }

  let potentialProfit;
  let potentialLoss;
  let kellyPercen;
  $: if (signals && signals.length > 0 && candles && candles.length > 1) {
    let wins = 0;
    let losses = 0;
    let totalProfit = 0;
    let totalLoss = 0;
    signals.forEach((s) => {
      if (s.Signal === 1) {
        if (s.Price < weightedPeriodAvg) {
          let profit = weightedPeriodAvg / s.Price - 1;
          wins += 1;
          totalProfit += profit;
        }
        if (s.Price > weightedPeriodAvg) {
          let loss = s.Price / weightedPeriodAvg - 1;
          losses += 1;
          totalLoss += loss;
        }
      }
      // else if (s.Signal === 2) {
      //   if (s.Price > weightedPeriodAvg) {
      //     let profit = s.Price / weightedPeriodAvg - 1;
      //     wins += 1;
      //     totalProfit += profit;
      //   }
      //   if (s.Price < weightedPeriodAvg) {
      //     let loss = weightedPeriodAvg / s.Price - 1;
      //     losses += 1;
      //     totalLoss += loss;
      //   }
      // }
    });
    potentialProfit = ((totalProfit / wins) * 100).toFixed(2);
    potentialLoss = ((totalLoss / losses) * 100).toFixed(2);
    kellyPercen = (
      (+winningRate / +potentialLoss -
        (100 - +winningRate) / +potentialProfit) *
      100
    ).toFixed(2);
  }
</script>

<div class="px-4 pt-4">
  <Autocomplete
    minCharactersToSearch={2}
    keywordsFunction={(it) => (it.symbol ? `${it.symbol.toLowerCase()}|^|${it.text.toLowerCase()}` : '')}
    on:change={stockChanged}
    outlined
    bind:value={stock}
    labelFieldName="text"
    label="Enter Company Name"
    items={$instruments} />
</div>
{#await loadChartPromise}
  <div class="px-4">
    Analyzing...
    <Spinner height="h-10" width="h-10" color="text-orange-600" />
  </div>
{:then _}
  <div class={candleClass}>
    {#if showAnalyzing}
      <div>
        <Progress fillColor="bg-orange-600" trackColor="bg-orange-200" />
      </div>
    {/if}
    <CandleChart {candles} {signals} />

    <div class="md:mb-4 flex flex-wrap">
      {#if candles && candles.length > 1}
        <div class="w-full md:w-1/3 px-0 md:pr-2 mb-0">
          <Input readonly value={increase} label="Last day / first day (%)" />
        </div>
      {/if}
      {#if signals && signals.length > 0}
        <div class="w-full md:w-1/3 px-0 mb-0">
          <Input readonly value={buyIncrease} label="Buy Increase (%)" />
        </div>
        <div class="w-full md:w-1/3 px-0 md:pl-2 mb-0">
          <Input
            readonly
            value={buySellIncrease}
            label="Buy+Sell Increase (%)" />
        </div>
      {/if}
    </div>

    <div class="md:mb-4 flex flex-wrap">
      {#if candles && candles.length > 1}
        <div class="w-full md:w-1/5 px-0 md:pr-2 mb-0">
          <Input
            readonly
            value={weightedPeriodAvg}
            label="Weighted Average Price" />
        </div>
      {/if}
      {#if signals && signals.length > 0}
        <div class="w-full md:w-1/5 px-0 mb-0">
          <Input readonly value={winningRate} label="Winning Rate (%)" />
        </div>
        <div class="w-full md:w-1/5 px-0 md:pl-2 mb-0">
          <Input readonly value={potentialProfit} label="Potential Profit" />
        </div>
        <div class="w-full md:w-1/5 px-0 md:pl-2 mb-0">
          <Input readonly value={potentialLoss} label="Potential Loss" />
        </div>
        <div class="w-full md:w-1/5 px-0 md:pl-2 mb-0">
          <Input readonly value={kellyPercen} label="Kelly %" />
        </div>
      {/if}
    </div>

    <div class="mb-5">
      <div>{periodLabel}</div>
      <Slider bind:value={freq} min={15} max={40} step={1} />
    </div>
    <div class="flex flex-wrap mb-5">
      <div class="w-full md:w-1/2 px-0 md:pr-2 mb-5 md:mb-0">
        <div>{buyFreqLabel}</div>
        <Slider bind:value={buyFreq} min={0} max={100} step={1} />
      </div>
      <div class="w-full md:w-1/2 px-0 md:pl-2">
        <div>{sellFreqLabel}</div>
        <Slider bind:value={sellFreq} min={0} max={100} step={1} />
      </div>
    </div>
    <div class="my-2">
      {#if $loginUser}
        <Button
          lg
          block
          bgColor="bg-orange-600"
          textColor="text-white"
          on:click={subcribe}>
          Subscribe
        </Button>
      {:else}
        <Button
          lg
          block
          bgColor="bg-orange-600"
          textColor="text-white"
          on:click={() => ($showSignIn = true)}>
          Login to subscribe
        </Button>
      {/if}
    </div>
  </div>
{/await}
