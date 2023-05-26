<route lang="yaml">
path: "/wallets"
</route>

<template>
  <main class="column bg-color px-5 py-5" id="main">
    <div id="title" class="sub-title">Wallets</div>

    <div class="table-container box px-0 py-2">
      <div class="px-1" style="padding-bottom: 0.5rem">
        <span>ALL ({{ wallets.length }})</span>
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
          <button class="button is-link is-outlined" @click="modal.open">
            <span>Add Wallet</span>
            <span class="icon is-small"><i class="fa-solid fa-plus"></i></span>
          </button>
          <button class="button is-white" title="刷新" @click="getData">
            <span class="icon is-small i-color"
              ><i class="fa-solid fa-arrows-rotate"></i
            ></span>
          </button>
          <button
            class="button is-white"
            title="全屏"
            @click="maximize"
            id="maximize"
          >
            <span class="icon is-small i-color"
              ><i class="fa-solid fa-maximize"></i
            ></span>
          </button>
          <button
            class="button is-white is-hidden"
            title="退出全屏"
            @click="minimize"
            id="minimize"
          >
            <span class="icon is-small i-color"
              ><i class="fa-solid fa-minimize"></i
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
            <th class="th-color">Address</th>
            <th class="th-color">Nodes</th>
            <th class="th-color">Balance</th>
            <th class="th-color">1 Day</th>
            <th class="th-color">Group</th>
          </tr>
        </thead>
        <tbody style="color: #757981; font-size: 14px">
          <tr v-for="w in wallets" :key="w.name">
            <td>
              <input
                type="checkbox"
                class="ckbox"
                :value="w.address"
                @click="checkOne"
              />
            </td>
            <td>
              <a :href="formatAddress(w.address)" target="_blank">{{
                w.name
              }}</a>
            </td>
            <td>{{ w.address }}</td>
            <td>
              {{ w.nodes[0] }}
              <span v-if="w.nodes[1] && w.nodes[1] > 0"
                >/ <span class="has-text-danger">{{ w.nodes[1] }}</span></span
              >
            </td>
            <td>{{ w.balance }} FIL</td>
            <td class="ok">+ {{ w.daily }} FIL</td>
            <td>{{ w.group }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <WalletModal ref="modal" />
  </main>
</template>

<script setup>
const route = useRoute();
const modal = ref();
const wallets = ref([]);
const origins = ref([]);
const keyword = ref("");
const api = useApi();
const isHidden = ref(true);
let addrs = [];

const checkAll = (e) => {
  let $this = e.target;
  const ckboxs = document.querySelectorAll(".ckbox");
  if ($this.checked) {
    ckboxs.forEach((e) => {
      e.checked = true;
      if (!addrs.includes(e.value)) {
        addrs.push(e.value);
      }
    });
    showDel();
  } else {
    ckboxs.forEach((e) => {
      e.checked = false;
    });
    showDel(false);
    addrs = [];
  }
};

const checkOne = (e) => {
  let $this = e.target;
  if ($this.checked) {
    showDel();
    if (!addrs.includes($this.value)) {
      addrs.push($this.value);
    }
  } else {
    const index = addrs.indexOf($this.value);
    if (index > -1) {
      addrs.splice(index, 1);
    }
    if (addrs.length === 0) {
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
      .delete("/api/wallets", {
        params: {
          addrs: addrs.toString(),
        },
      })
      .then(() => {
        addrs = [];
        showDel(false);
        getData();
      });
  }
};

const formatAddress = (addr) => {
  return "https://dashboard.saturn.tech/address/" + addr + "?period=1+Days";
};

const maximize = () => {
  document.getElementById("navbar").classList.add("is-hidden");
  document.getElementById("aside").classList.add("is-hidden");
  document.getElementById("title").classList.add("is-hidden");
  document.getElementById("maximize").classList.add("is-hidden");
  document.getElementById("minimize").classList.remove("is-hidden");
};

const minimize = () => {
  document.getElementById("navbar").classList.remove("is-hidden");
  document.getElementById("aside").classList.remove("is-hidden");
  document.getElementById("title").classList.remove("is-hidden");
  document.getElementById("maximize").classList.remove("is-hidden");
  document.getElementById("minimize").classList.add("is-hidden");
};

const getData = () => {
  api.get("/api/wallets").then((res) => {
    wallets.value = res.data;
    origins.value = res.data;
  });
};

let timeoutId;
const onKeyUp = () => {
  clearTimeout(timeoutId); // 每次先清除已有的定时器
  timeoutId = setTimeout(() => {
    // 这里编写要执行的操作代码...
    let key = keyword.value;
    if (key) {
      key = key.toLowerCase();
      wallets.value = origins.value.filter((o) => {
        return (
          o.name.toLowerCase().includes(key) ||
          o.address.includes(key) ||
          (o.group && o.group.toLowerCase().includes(key))
        );
      });
    } else {
      wallets.value = origins.value;
    }
  }, 500); // 设置延时
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
