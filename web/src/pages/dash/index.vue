<template>
  <main class="column bg-color px-5 py-5" id="main">
    <div style="margin-bottom: 0.8rem">Summary</div>
    <div class="columns">
      <div class="column is-3">
        <div class="card">
          <header class="card-header">
            <p class="card-header-title">SATURN STATUS</p>
            <span class="card-header-icon" aria-label="more options">
              <span class="icon">
                <i
                  class="fa-brands fa-connectdevelop has-text-primary"
                  aria-hidden="true"
                ></i>
              </span>
            </span>
          </header>
          <div class="card-image">
            <figure class="image mt-1">
              <apexchart
                type="donut"
                :options="options"
                :series="series"
              ></apexchart>
            </figure>
          </div>
          <div class="card-content">
            <nav class="level">
              <div class="level-item has-text-centered">
                <div>
                  <p class="title">
                    {{ nodes[0] }}
                    <span v-if="nodes[1] && nodes[1] > 0"
                      >/
                      <span class="has-text-danger">{{ nodes[1] }}</span></span
                    >
                  </p>
                  <p class="heading">My Nodes</p>
                </div>
              </div>
              <div class="level-item has-text-centered">
                <div>
                  <p class="title">{{ earnings }}</p>
                  <p class="heading">Wallet FIL</p>
                </div>
              </div>
            </nav>
          </div>
        </div>
      </div>
      <div class="column">
        <div class="box">
          <MonthChart :data="dailys" :time="time"></MonthChart>
        </div>
        <div class="box">
          <GroupChart :data="groups"></GroupChart>
        </div>
      </div>
    </div>
  </main>
</template>

<script setup>
import { ref } from "vue";

const api = useApi();
const earnings = ref(0);
const nodes = ref([]);
const groups = ref([]);
const time = ref("");
const dailys = ref([]);

const options = {
  legend: {
    position: "bottom",
  },
  colors: ["#00E396", "#FEB019"],
  labels: ["Active", "Inactive"],
};
const series = ref([]);

api.get("/api/summary").then((res) => {
  nodes.value = res.data.nodes;
  series.value = [res.data.nodes[0], res.data.nodes[1]];

  time.value = res.data.time;
  dailys.value = res.data.dailys;

  groups.value = res.data.groups;
});
</script>
