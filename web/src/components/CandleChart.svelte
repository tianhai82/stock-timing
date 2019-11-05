<script>
  import echarts from "echarts";
  import ECharts from "echarts-for-svelte";

  export let candles = [];
  export let signals = [];
  let buys = [];
  let sells = [];
  let dates = [];
  let data = [];
  let option = {};

  $: {
    dates = candles.map(c => c.FromDate.substr(0, 10));
    data = candles.map(c => [c.Open, c.Close, c.Low, c.High]);

    buys = signals
      .filter(s => s.Signal === 1)
      .map(s => [s.Date.substr(0, 10), s.Price]);
    sells = signals
      .filter(s => s.Signal === 2)
      .map(s => [s.Date.substr(0, 10), s.Price]);

    const candleSeries = { ...optionTemplate.series[0], data };
    const buySeries = { ...optionTemplate.series[1], data: buys };
    const sellSeries = { ...optionTemplate.series[2], data: sells };
    option = {
      ...optionTemplate,
      xAxis: { ...optionTemplate.xAxis, data: dates },
      series: [candleSeries, buySeries, sellSeries]
    };
  }

  const colorList = [
    "#c23531",
    "#2f4554",
    "#61a0a8",
    "#d48265",
    "#91c7ae",
    "#749f83",
    "#ca8622",
    "#bda29a",
    "#6e7074",
    "#546570",
    "#c4ccd3"
  ];

  const optionTemplate = {
    animation: false,
    color: colorList,
    tooltip: {
      triggerOn: "none",
      transitionDuration: 0,
      confine: true,
      borderRadius: 4,
      borderWidth: 1,
      borderColor: "#333",
      backgroundColor: "rgba(255,255,255,0.9)",
      textStyle: {
        fontSize: 12,
        color: "#333"
      },
      position: function(pos, params, el, elRect, size) {
        var obj = {
          top: 60
        };
        obj[["left", "right"][+(pos[0] < size.viewSize[0] / 2)]] = 10;
        return obj;
      }
    },
    dataZoom: [
      {
        type: "slider",
        realtime: false,
        start: 40,
        end: 100,
        top: 5,
        height: 20,
        handleIcon:
          "M10.7,11.9H9.3c-4.9,0.3-8.8,4.4-8.8,9.4c0,5,3.9,9.1,8.8,9.4h1.3c4.9-0.3,8.8-4.4,8.8-9.4C19.5,16.3,15.6,12.2,10.7,11.9z M13.3,24.4H6.7V23h6.6V24.4z M13.3,19.6H6.7v-1.4h6.6V19.6z",
        handleSize: "120%"
      },
      {
        type: "inside",
        start: 50,
        end: 70,
        top: 30,
        height: 20
      }
    ],
    xAxis: {
      type: "category",
      splitNumber: 5,
      data: [],
      boundaryGap: false,
      axisLine: { lineStyle: { color: "#777" } },
      axisLabel: {
        formatter: function(value) {
          return echarts.format.formatTime("MM-dd", value);
        }
      },
      min: "dataMin",
      max: "dataMax",
      axisPointer: {
        label: { show: false },
        triggerTooltip: true,
        handle: {
          show: true,
          margin: 30,
          color: "#B80C00"
        }
      }
    },
    yAxis: {
      scale: true,
      splitNumber: 5,
      axisLine: { lineStyle: { color: "#777" } },
      splitLine: { show: true },
      axisTick: { show: false },
      axisLabel: {
        inside: true,
        formatter: "{value}\n"
      }
    },
    grid: {
      left: 2,
      right: 5,
      top: 30,
      height: 315
    },
    graphic: {
      type: "group",
      left: "center",
      top: 70,
      width: 300,
      bounding: "raw"
    },
    series: [
      {
        type: "candlestick",
        name: "Price",
        data: [],
        itemStyle: {
          normal: {
            color: "#14b143",
            color0: "#ef232a",
            borderColor: "#14b143",
            borderColor0: "#ef232a"
          },
          emphasis: {
            color: "black",
            color0: "#444",
            borderColor: "black",
            borderColor0: "#444"
          }
        }
      },
      {
        type: "scatter",
        name: "Buys",
        symbolSize: 20,
        symbol: "arrow",
        itemStyle: {
          color: "#5ccc7e"
        },
        data: buys
      },
      {
        type: "scatter",
        name: "Sells",
        symbolSize: 20,
        symbolRotate: 180,
        symbol: "arrow",
        itemStyle: {
          color: "#ff8673"
        },
        data: sells
      }
    ]
  };
</script>

<ECharts {echarts} {option} style="height:400px" />
