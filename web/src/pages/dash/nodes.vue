<route lang="yaml">
path: "/nodes"
</route>
<template>
  <main class="column bg-color px-5 py-5" id="main">
    <div class="sub-title">Nodes</div>

    <div class="table-container box px-0 py-2">
      <div class="px-1" style="padding-bottom: 0.5rem">
        <span>ALL ({{ nodes.length }})</span>
        <input
          class="input is-rounded is-small ml-1"
          type="text"
          placeholder="Filter"
          style="max-width: 20%"
          v-model="keyword"
          @keyup="onKeyUp"
        />

        <div class="buttons are-small is-pulled-right">
          <button
            class="button is-danger is-outlined"
            :class="{ 'is-hidden': isHidden }"
            title="删除"
            @click="del"
          >
            <span>Delete</span>
            <span class="icon is-small"
              ><i class="fa-solid fa-trash-can"></i
            ></span>
          </button>
          <button
            class="button is-warning is-outlined"
            :class="{ 'is-hidden': isHidden }"
            @click="upgrade"
          >
            <span>Upgrade</span>
            <span class="icon is-small"
              ><i class="fa-brands fa-docker"></i
            ></span>
          </button>
          <button
            class="button is-info is-outlined"
            :class="{ 'is-hidden': isHidden }"
            @click="poolModal.open()"
          >
            <span>Traffic Pool</span>
            <span class="icon is-small"><i class="fa-solid fa-water"></i></span>
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
            <th>
              <input
                type="checkbox"
                class="check"
                @click="checkAll"
                id="checkAll"
              />
            </th>
            <th class="th-color">Name</th>
            <th class="th-color">State</th>
            <th class="th-color">ID</th>
            <th class="th-color">Type</th>
            <th class="th-color">IP</th>
            <th class="th-color">CPU</th>
            <th class="th-color">RAM</th>
            <th class="th-color">Bandwidth</th>
            <th class="th-color">Disk</th>
            <th class="th-color">Traffic</th>
            <th class="th-color">Price</th>
            <th class="th-color">Renew</th>
            <th class="th-color">
              <i class="fa-solid fa-ellipsis-vertical"></i>
            </th>
          </tr>
        </thead>
        <tbody style="color: #757981; font-size: 14px">
          <tr v-for="e in nodes" :key="e.id">
            <td>
              <input
                type="checkbox"
                class="ckbox"
                :value="e.ip"
                @click="checkOne"
              />
            </td>
            <td>
              <a>{{ e.name }}</a>
            </td>
            <td
              :class="
                e.state === 'active'
                  ? 'has-text-success'
                  : e.state === 'draining'
                  ? 'has-text-warning'
                  : 'has-text-danger'
              "
            >
              {{ e.state }}
            </td>
            <td>{{ shortNodeId(e.nodeId) }}</td>
            <td>{{ e.type }}</td>
            <td>{{ e.ip }}</td>
            <td>{{ e.cpu }}</td>
            <td>{{ e.ram }}</td>
            <td>{{ e.bandwidth }}Gbps</td>
            <td style="white-space: pre">{{ prettifyDisk(e.disk) }}</td>
            <td>{{ e.traffic }}</td>
            <td>${{ e.price }}</td>
            <td>{{ e.renew }}</td>
            <td class="i-color pointer" @click="terminal.open(e.name, e.ip)">
              <i class="fa-solid fa-terminal"></i>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <NodeModal ref="nodeModal" />
    <PoolModal ref="poolModal" :ips="ips" />
    <TerminalModal ref="terminal" />
  </main>
</template>

<script setup>
const route = useRoute();
const api = useApi();
let ips = [];
const isHidden = ref(true);
const nodeModal = ref();
const poolModal = ref();
const terminal = ref();

const nodes = ref([]);
const origins = ref([]);
const keyword = ref("");

const checkAll = (e) => {
  let $this = e.target;
  const ckboxs = document.querySelectorAll(".ckbox");
  if ($this.checked) {
    ckboxs.forEach((e) => {
      e.checked = true;
      if (!ips.includes(e.value)) {
        ips.push(e.value);
      }
    });
    showDel();
  } else {
    ckboxs.forEach((e) => {
      e.checked = false;
    });
    showDel(false);
    ips = [];
  }
};

const checkOne = (e) => {
  let $this = e.target;
  if ($this.checked) {
    showDel();
    if (!ips.includes($this.value)) {
      ips.push($this.value);
    }
  } else {
    const index = ips.indexOf($this.value);
    if (index > -1) {
      ips.splice(index, 1);
    }
    if (ips.length === 0) {
      showDel(false);
      document.getElementById("checkAll").checked = false;
    }
  }
};

const showDel = (hide) => {
  isHidden.value = hide === false;
};

const del = () => {
  if (confirm("确定删除吗？")) {
    api
      .delete("/api/nodes", {
        params: {
          ip: ips.toString(),
        },
      })
      .then(() => {
        ips = [];
        showDel(false);
        getData();
      });
  }
};

const upgrade = () => {
  if (confirm("确定升级节点版本吗？")) {
    api
      .post(
        "/api/nodes/upgrade",
        { ip: ips.toString() },
        {
          headers: {
            "Content-Type": "application/x-www-form-urlencoded",
          },
        }
      )
      .then(() => {
        alert("升级节点版本请求已发送");
      });
  }
};

const getData = () => {
  api.get("/api/nodes").then((res) => {
    nodes.value = res.data;
    origins.value = res.data;
  });
};

let timeoutId;
const onKeyUp = () => {
  clearTimeout(timeoutId); // 每次先清除已有的定时器
  timeoutId = setTimeout(() => {
    let key = keyword.value;
    if (key) {
      key = key.toLowerCase();
      nodes.value = origins.value.filter((o) => {
        return (
          o.name.toLowerCase().includes(key) ||
          o.ip.includes(key) ||
          (o.group && o.group.toLowerCase().includes(key)) ||
          (o.state && o.state.toLowerCase().includes(key)) ||
          (o.nodeId && o.nodeId.toLowerCase().includes(key))
        );
      });
    } else {
      nodes.value = origins.value;
    }
  }, 500); // 设置延时
};

const shortNodeId = (e) => {
  if (e && e.length > 8) {
    const parts = e.split("-");
    return parts[0];
  }
  return e;
};

const prettifyDisk = (e) => {
  // 158G 11G 7%
  if (e && e.length > 7) {
    const parts = e.split(" ");
    if (parts.length < 3) {
      return e;
    }
    let size = parts[0];
    let used = parts[1];
    let percent = "(" + parts[2] + ")";
    if (used.length < 4) {
      used = used.padStart(4, " ");
    }
    if (size.length < 4) {
      size = size.padStart(4, " ");
    }
    if (percent.length < 5) {
      percent = percent.padStart(5, " ");
    }
    return used + "/" + size + " " + percent;
  }

  return e;
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
