<template>
  <main class="column bg-color px-5 py-5" id="main">
    <div style="margin-bottom: 0.8rem;">Summary</div>
    <div class="columns">
      <div class="column is-3">
        <div class="card">
            <header class="card-header">
              <p class="card-header-title">
                SATURN STATUS
              </p>
              <span class="card-header-icon" aria-label="more options">
                <span class="icon">
                  <i class="fa-brands fa-connectdevelop has-text-primary" aria-hidden="true"></i>
                </span>
              </span>
            </header>
            <div class="card-image">
              <figure class="image mt-1">
                <apexchart type="donut" :options="options" :series="series"></apexchart>
              </figure>
            </div>
            <div class="card-content">
              <nav class="level">
                <div class="level-item has-text-centered">
                  <div>
                    <p class="title">{{ nodeCount }}</p>
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
          <apexchart type="line" height="200" :options="monthOptions" :series="monthSeries"></apexchart>
        </div>
      </div>
    </div>
  </main>  
  </template>
  
<script setup>
  const earnings = ref(0)
  const nodeCount = ref(0)
  const api = useApi()

  const options = {
    legend: {
      position: 'bottom'
    },
    colors: ['#00E396', '#FEB019'],
    labels: ['Active', 'Inactive']
  }
  const series = ref([])

  const monthOptions = ref({})

  const monthSeries = ref([
    {
      name: "FIL",
      data: [],
    }
  ])

  api.get('/api/summary').then(res => {
    earnings.value = res.data.earnings
    nodeCount.value = res.data.nodes[0] + res.data.nodes[1]
    series.value = [res.data.nodes[0], res.data.nodes[1]]
    monthSeries.value[0].data = res.data.dailys
    monthOptions.value = {
      title: {
          text: 'Earnings by Month ' + res.data.time,
          align: 'left'
        },
        dataLabels: {
            enabled: false
        },
        chart: {
          toolbar: {
            show: false
        }
      }
    
    }
  })
</script>

  