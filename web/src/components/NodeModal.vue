<template>
  <div class="modal" id="AddNodeModal" :class="{ 'is-active': isOpen }">
    <div class="modal-background"></div>
    <div class="modal-card">
      <header class="modal-card-head">
        <p class="modal-card-title">Add Node</p>
        <button class="delete" aria-label="close" @click="close"></button>
      </header>
      <section class="modal-card-body">
        <div class="field is-horizontal">
          <div class="field-label is-normal">
            <label class="label"
              ><span class="has-text-danger">*</span> Name</label
            >
          </div>
          <div class="field-body">
            <div class="field">
              <div class="control">
                <input
                  class="input"
                  type="text"
                  maxlength="24"
                  placeholder="Name"
                  v-model.lazy.trim="name"
                />
              </div>
            </div>
          </div>
        </div>
        <div class="field is-horizontal">
          <div class="field-label is-normal">
            <label class="label"
              ><span class="has-text-danger">*</span> IP</label
            >
          </div>
          <div class="field-body">
            <div class="field">
              <div class="control">
                <input
                  class="input"
                  type="text"
                  maxlength="15"
                  placeholder="IP"
                  v-model.lazy.trim="ip"
                />
              </div>
            </div>
          </div>
        </div>
        <div class="field is-horizontal">
          <div class="field-label is-normal">
            <label class="label">Bandwidth</label>
          </div>
          <div class="field-body">
            <div class="field has-addons">
              <p class="control">
                <input
                  class="input"
                  type="number"
                  maxlength="32"
                  placeholder="Bandwidth"
                  v-model.lazy.trim="bandwidth"
                />
              </p>
              <p class="control">
                <a class="button">Gbps</a>
              </p>
            </div>
          </div>
        </div>
        <div class="field is-horizontal">
          <div class="field-label is-normal">
            <label class="label">Traffic</label>
          </div>
          <div class="field-body">
            <div class="field has-addons">
              <div class="control">
                <input
                  class="input"
                  type="number"
                  maxlength="32"
                  placeholder="Traffic"
                  v-model.lazy.trim="traffic"
                />
              </div>
              <p class="control">
                <a class="button">TB</a>
              </p>
            </div>
          </div>
        </div>
        <div class="field is-horizontal">
          <div class="field-label is-normal">
            <label class="label">Price</label>
          </div>
          <div class="field-body">
            <div class="field has-addons">
              <div class="control">
                <input
                  class="input"
                  type="number"
                  maxlength="32"
                  placeholder="Price"
                  v-model.lazy.trim="price"
                />
              </div>
              <p class="control">
                <a class="button">$</a>
              </p>
            </div>
          </div>
        </div>
        <div class="field is-horizontal">
          <div class="field-label is-normal">
            <label class="label">Renew</label>
          </div>
          <div class="field-body">
            <div class="field">
              <div class="control">
                <input
                  class="input"
                  type="date"
                  maxlength="32"
                  placeholder="Renew"
                  v-model.lazy.trim="renew"
                />
              </div>
            </div>
          </div>
        </div>
      </section>
      <footer class="modal-card-foot">
        <button class="button is-link is-fullwidth" @click="add">Add</button>
      </footer>
    </div>
  </div>
</template>

<script setup>
const api = useApi();
const isOpen = ref(false);
const name = ref("");
const ip = ref("");
const bandwidth = ref();
const traffic = ref();
const price = ref();
const renew = ref();

const open = () => {
  isOpen.value = true;
};

const close = () => {
  isOpen.value = false;
};

const add = () => {
  if (!name.value || !ip.value) {
    alert("名字和ip不能为空");
    return;
  }
  let trafficVal = 0;
  if (traffic.value) {
    trafficVal = traffic.value;
  }
  let bandwidthVal = 0;
  if (bandwidth.value) {
    bandwidthVal = bandwidth.value;
  }
  let priceVal = 0;
  if (price.value) {
    priceVal = price.value;
  }
  let renewVal = "-";
  if (renew.value && renew.value.length > 0) {
    renewVal = renew.value;
  }
  api
    .post("/api/nodes", {
      name: name.value,
      ip: ip.value,
      bandwidth: bandwidthVal,
      traffic: trafficVal + "TB",
      price: priceVal,
      renew: renewVal,
    })
    .then(() => {
      ip.value = "";
    });
};

defineExpose({ open });
</script>
