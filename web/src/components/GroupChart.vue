<template>
  <div id="groupChart">
    <apexchart
      type="bar"
      height="220"
      :options="chartOptions"
      :series="series"
    ></apexchart>
  </div>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
  data: Array,
});

const { data: groups } = toRefs(props);

const series = computed(() => {
  return [
    {
      name: "FIL",
      data: groups.value.map((e) => e.balance),
    },
  ];
});

const chartOptions = computed(() => {
  return {
    title: {
      text: "Group Earnings",
      align: "left",
    },
    xaxis: {
      categories: groups.value.map((e) => e.name),
      position: "bottom",
      axisBorder: {
        show: false,
      },
      axisTicks: {
        show: false,
      },
      tooltip: {
        enabled: false,
      },
    },
    chart: {
      toolbar: {
        show: false,
      },
      type: "bar",
    },
    plotOptions: {
      bar: {
        horizontal: false,
        columnWidth: "10%",
        borderRadius: 8,
        dataLabels: {
          position: "top", // top, center, bottom
        },
      },
    },
    dataLabels: {
      enabled: true,
      offsetY: -20,
      style: {
        //fontSize: "12px",
        colors: ["#304758"],
      },
    },

    yaxis: {
      axisBorder: {
        show: false,
      },
      axisTicks: {
        show: false,
      },
      labels: {
        show: false,
      },
    },
  };
});
</script>
