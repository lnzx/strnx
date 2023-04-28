<route lang="yaml">
    path: "/nodes"
</route>
<template>
    <main class="column bg-color px-5 py-5" id="main">
      <div class="sub-title">Nodes</div>
    
      <div class="table-container box px-0 py-2">
        <div class="px-1" style="padding-bottom: 0.5rem;">
          <span>ALL (5)</span>
    
          <div class="buttons are-small is-pulled-right">
            <button class="button is-danger is-outlined is-hidden">
                <span>Delete</span>
                <span class="icon is-small"><i class="fa-solid fa-trash-can"></i></span>
            </button>
            <button class="button is-link is-outlined" @click="updateNode">
                <span>Update Node</span>
                <span class="icon is-small"><i class="fa-solid fa-plus"></i></span>
            </button>
            <button class="button is-link is-outlined" @click="openModal">
                <span>Add Node</span>
                <span class="icon is-small"><i class="fa-solid fa-plus"></i></span>
            </button>
            <button class="button is-white" title="refresh">
                <span class="icon is-small i-color"><i class="fa-solid fa-arrows-rotate"></i></span>
            </button>
          </div>
        </div>
        <table class="table is-fullwidth">
            <thead style="background: #f8f8fb; font-size:14px;">
                <tr>
                    <th><input type="checkbox" class="check" @click="checkAll"></th>
                    <th class="th-color">ID</th>
                    <th class="th-color">IP</th>
                    <th class="th-color">ISP</th>
                    <th class="th-color">Location</th>
                    <th class="th-color">CPU</th>
                    <th class="th-color">RAM</th>
                    <th class="th-color">Disk</th>
                    <th class="th-color">Bandwidth</th>
                    <th class="th-color">Traffic</th>
                    <th class="th-color">Price</th>
                    <th class="th-color">Renew</th>
                </tr>
            </thead>
            <tbody style="color:#757981; font-size:14px;">
                <tr>
                    <td><input type="checkbox" class="ckbox"></td>
                    <td><a>423667f8</a></td>
                    <td>44.234.105.210</td>
                    <td>Digital Ocean</td>
                    <td>Singapore (SG)</td>
                    <td>6</td>
                    <td>32</td>
                    <td>1TB</td>
                    <td>10Gbps</td>
                    <td>unlimit</td>
                    <td>$300</td>
                    <td>2022/12/20</td>
                </tr>
            </tbody>
        </table>
      </div>

      <NodeModal ref="nodeModal" />
    </main>  
</template>
    
<script setup>
    const route = useRoute()
    const api = useApi()
    const nodeModal = ref()
    const id = ref('')
    const ip = ref('')
    const nodes = ref([])

    const checkAll = (e) => {
        let $this = e.target;
        const ckboxs = document.querySelectorAll('.ckbox')
        if($this.checked){
            ckboxs.forEach(e => {
                e.checked = true
            })
        }else{
            ckboxs.forEach(e => {
                e.checked = false
            })
        }
    }

    const openModal = () => {
        console.log('click open modal')
        nodeModal.value.open()
    }

    const updateNode = () => {
        console.log('updateNode')
    }
    
    const getData = (node) => {
        api.get('/api/nodes').then(res => {
            let data = res.data
            if(node){
                let filters = []
                data.forEach(e => {
                    if(e.address == node){
                        filters.push(e)
                    }
                })
                nodes.value = filters
            }else{
                nodes.value = data
            }
        })
    }

    getData(route.query.node)
</script>
    
<style scoped>
    .warn {color:#f4a271;}
    .ok {color: #48c78e;}
    .err {color: #f14668;}
    .check {color:rgba(59,130,246,0.5);}
    .th-color {color: #54565b;}
    .pointer {cursor: pointer;}
    .i-color {color: #183153;}
</style>
    
    
