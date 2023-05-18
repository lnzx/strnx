<route lang="yaml">
path: "/nodes"
</route>
<template>
  <main class="column bg-color px-5 py-5" id="main">
    <div class="sub-title">Nodes</div>

    <div class="table-container box px-0 py-2">
      <div class="px-1" style="padding-bottom: 0.5rem">
        <span>ALL ({{ nodes.length }})</span>

        <div class="buttons are-small is-pulled-right">
          <button class="button is-danger is-outlined" :class="{ 'is-hidden': isHidden }" title="删除" @click="del">
            <span>Delete</span>
            <span class="icon is-small"><i class="fa-solid fa-trash-can"></i></span>
          </button>
          <button class="button is-warning is-outlined" :class="{ 'is-hidden': isHidden }" @click="upgrade">
            <span>Upgrade</span>
            <span class="icon is-small"><i class="fa-brands fa-docker"></i></span>
          </button>
          <button class="button is-link is-outlined" @click="nodeModal.open()">
            <span>Add Node</span>
            <span class="icon is-small"><i class="fa-solid fa-plus"></i></span>
          </button>
          <button class="button is-white" title="refresh" @click="getData">
            <span class="icon is-small i-color"
              ><i class="fa-solid fa-arrows-rotate"></i
            ></span>
          </button>
        </div>
      </div>
      <table class="table is-fullwidth">
        <thead style="background: #f8f8fb; font-size: 14px">
          <tr>
            <th><input type="checkbox" class="check" @click="checkAll" id="checkAll"></th>
            <th class="th-color">Name</th>
            <th class="th-color">State</th>
            <th class="th-color">Type</th>
            <th class="th-color">IP</th>
            <th class="th-color">CPU</th>
            <th class="th-color">RAM</th>
            <th class="th-color">Bandwidth</th>
            <th class="th-color">Disk</th>
            <th class="th-color">Traffic</th>
            <th class="th-color">Price</th>
            <th class="th-color">Renew</th>
          </tr>
        </thead>
        <tbody style="color: #757981; font-size: 14px">
          <tr v-for="e in nodes" :key="e.id">
            <td><input type="checkbox" class="ckbox" :value="e.ip" @click="checkOne" /></td>
            <td><a>{{ e.name }}</a></td>
            <td class="has-text-success">{{ e.state }}</td>
            <td>{{ e.type }}</td>
            <td>{{ e.ip }}</td>
            <td>{{ e.cpu }}</td>
            <td>{{ e.ram }}</td>
            <td>{{ e.bandwidth }}Gbps</td>
            <td>{{ e.disk }}</td>
            <td>{{ e.traffic }}</td>
            <td>${{ e.price }}</td>
            <td>{{ e.renew }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <NodeModal ref="nodeModal" />
  </main>
</template>

<script setup>
const route = useRoute();
const api = useApi();
let ips = [];
const isHidden = ref(true)
const nodeModal = ref();
const nodes = ref([]);

const checkAll = (e) => {
  let $this = e.target;
  const ckboxs = document.querySelectorAll('.ckbox')
  if($this.checked){
    ckboxs.forEach(e => {
      e.checked = true
      if(!ips.includes(e.value)){
        ips.push(e.value)
      }
    })
    showDel()
  }else{
    ckboxs.forEach(e => {
      e.checked = false
    })
    showDel(false)
    ips = []
  }
}

const checkOne = (e) => {
  let $this = e.target;
  if($this.checked){
    showDel()
    if(!ips.includes($this.value)){
      ips.push($this.value)
    }
  }else{
    const index = ips.indexOf($this.value);
    if(index > -1){
      ips.splice(index, 1)
    }
    if(ips.length === 0){
      showDel(false)
      document.getElementById('checkAll').checked = false
    }
  }
}

const showDel = (hide) => {
  isHidden.value = (hide === false)
}

const del = () => {
  if(confirm('确定删除吗？')){
    api.delete('/api/nodes', {params: {
        'ip': ips.toString()
      }}).then(() => {
      ips = []
      showDel(false)
      getData()
    })
  }
}

const upgrade = () => {
  if(confirm('确定升级节点版本吗？')){
    api.post('/api/nodes/upgrade', {ip: ips.toString()}, {
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded'
      }
    }).then(() => {
      alert('升级节点版本请求已发送。')
    })
  }
};

const getData = () => {
  api.get("/api/nodes").then((res) => {
    nodes.value = res.data;
  });
};

getData();
</script>

<style scoped>
.check {
  color: rgba(59, 130, 246, 0.5);
}
.th-color {
  color: #54565b;
}
.i-color {
  color: #183153;
}
</style>
