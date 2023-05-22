<template>
  <div class="modal" id="AddNodeModal" :class="{ 'is-flex': isOpen }">
    <div class="modal-background"></div>
    <div class="modal-card">
      <header class="modal-card-head">
        <p class="modal-card-title">Create Traffic Pool</p>
        <button class="delete" aria-label="close" @click="close"></button>
      </header>
      <section class="modal-card-body">
        <div class="field is-horizontal">
          <div class="field-label is-normal">
            <label class="label"
              ><span class="has-text-danger">*</span> Traffic</label
            >
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
              <p class="control">
                <span
                  class="button is-success"
                  :class="{ 'is-hidden': isHideOk }"
                  >Success</span
                >
                <span
                  class="button is-danger"
                  :class="{ 'is-hidden': isHideErr }"
                  >流量不能为空</span
                >
              </p>
            </div>
          </div>
        </div>
      </section>
      <footer class="modal-card-foot">
        <button class="button is-link is-fullwidth" @click="add">Create</button>
      </footer>
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue";

const isOpen = ref(false);
const traffic = ref();
const api = useApi();
const isHideOk = ref(true);
const isHideErr = ref(true);

const props = defineProps({
  ips: Array,
});

const open = () => {
  isHideErr.value = true;
  isHideOk.value = true;
  isOpen.value = true;
};

const close = () => {
  isOpen.value = false;
};

const add = () => {
  if (!traffic.value) {
    isHideErr.value = false;
    isHideOk.value = true;
    return;
  }
  isHideErr.value = true;
  api
    .post(
      "/api/nodes/pool",
      {
        traffic: traffic.value,
        ip: props.ips.toString(),
      },
      {
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
        },
      }
    )
    .then(() => {
      traffic.value = "";
      isHideOk.value = false;
    });
};

defineExpose({ open });
</script>
