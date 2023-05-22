<template>
  <apexchart
    type="line"
    height="200"
    :options="options"
    :series="series"
  ></apexchart>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
  data: Array,
  time: String,
});

const { data: dailys, time } = toRefs(props);

const series = computed(() => {
  return [
    {
      name: "FIL",
      data: dailys.value,
    },
  ];
});

const options = computed(() => {
  return {
    title: {
      text: "Earnings by Month " + time.value,
      align: "left",
    },
    dataLabels: {
      enabled: false,
    },
    chart: {
      toolbar: {
        show: false,
      },
    },
  };
});
</script>
