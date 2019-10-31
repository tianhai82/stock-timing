<script>
  import echarts from "echarts";
  import ECharts from "echarts-for-svelte";

  export let candles = [];
  let dates = [];
  let data = [];
  let option = {};
  $: {
    dates = candles.map(c => c.FromDate.substr(0, 10));
    data = candles.map(c => [c.Open, c.Close, c.Low, c.High]);
    option = {
      ...optionTemplate,
      xAxis: [
        { ...optionTemplate.xAxis[0], data: dates },
        { ...optionTemplate.xAxis[1], data: dates }
      ],
      series: [{ ...optionTemplate.series[0], data }]
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

  let optionTemplate = {
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
        obj[["left", "right"][+(pos[0] < size.viewSize[0] / 2)]] = 5;
        return obj;
      }
    },
    dataZoom: [
      {
        type: "slider",
        // xAxisIndex: [0, 1],
        realtime: false,
        start: 50,
        end: 100,
        top: 15,
        height: 20,
        handleIcon:
          "M10.7,11.9H9.3c-4.9,0.3-8.8,4.4-8.8,9.4c0,5,3.9,9.1,8.8,9.4h1.3c4.9-0.3,8.8-4.4,8.8-9.4C19.5,16.3,15.6,12.2,10.7,11.9z M13.3,24.4H6.7V23h6.6V24.4z M13.3,19.6H6.7v-1.4h6.6V19.6z",
        handleSize: "120%"
      },
      {
        type: "inside",
        // xAxisIndex: [0, 1],
        start: 50,
        end: 70,
        top: 30,
        height: 20
      }
    ],
    xAxis: [
      {
        type: "category",
        data: dates,
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
          type: "shadow",
          label: { show: false },
          triggerTooltip: true,
          handle: {
            show: true,
            margin: 30,
            color: "#B80C00"
          }
        }
      }
      // {
      //   type: "category",
      //   gridIndex: 1,
      //   data: dates,
      //   scale: true,
      //   boundaryGap: false,
      //   splitLine: { show: false },
      //   axisLabel: { show: false },
      //   axisTick: { show: false },
      //   axisLine: { lineStyle: { color: "#777" } },
      //   splitNumber: 20,
      //   min: "dataMin",
      //   max: "dataMax",
      //   axisPointer: {
      //     type: "shadow",
      //     label: { show: false },
      //     triggerTooltip: true,
      //     handle: {
      //       show: true,
      //       margin: 30,
      //       color: "#B80C00"
      //     }
      //   }
      // }
    ],
    yAxis: [
      {
        scale: true,
        splitNumber: 2,
        axisLine: { lineStyle: { color: "#777" } },
        splitLine: { show: true },
        axisTick: { show: false },
        axisLabel: {
          inside: true,
          formatter: "{value}\n"
        }
      }
      // {
      //   scale: true,
      //   gridIndex: 1,
      //   splitNumber: 2,
      //   axisLabel: { show: false },
      //   axisLine: { show: false },
      //   axisTick: { show: false },
      //   splitLine: { show: false }
      // }
    ],
    grid: [
      {
        left: 0,
        right: 20,
        top: 10,
        height: 260
      }
      // {
      //   left: 20,
      //   right: 20,
      //   height: 40,
      //   top: 260
      // }
    ],
    graphic: [
      {
        type: "group",
        left: "center",
        top: 70,
        width: 300,
        bounding: "raw"
      }
    ],
    series: [
      {
        type: "candlestick",
        name: "Price",
        data: data,
        itemStyle: {
          normal: {
            color: "#ef232a",
            color0: "#14b143",
            borderColor: "#ef232a",
            borderColor0: "#14b143"
          },
          emphasis: {
            color: "black",
            color0: "#444",
            borderColor: "black",
            borderColor0: "#444"
          }
        }
      }
    ]
  };
</script>

<ECharts {echarts} {option} style="height:330px" />
