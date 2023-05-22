<template>
  <div class="modal" id="AddNodeModal" :class="{ 'is-flex': isOpen }">
    <div class="modal-background"></div>
    <div class="modal-card">
      <header class="modal-card-head">
        <p class="modal-card-title">Add Wallet</p>
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
              ><span class="has-text-danger">*</span> Address</label
            >
          </div>
          <div class="field-body">
            <div class="field">
              <div class="control">
                <input
                  class="input"
                  type="text"
                  maxlength="41"
                  placeholder="Address"
                  v-model.lazy.trim="address"
                />
              </div>
            </div>
          </div>
        </div>
        <div class="field is-horizontal">
          <div class="field-label is-normal">
            <label class="label">Group</label>
          </div>
          <div class="field-body">
            <div class="field">
              <div class="control">
                <input
                  class="input"
                  type="text"
                  maxlength="8"
                  placeholder="Group"
                  v-model.lazy.trim="group"
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
const isOpen = ref(false);
const name = ref("");
const address = ref("");
const group = ref("");
const api = useApi();

const open = () => {
  isOpen.value = true;
};

const close = () => {
  isOpen.value = false;
};

const add = () => {
  if (!name.value || !address.value || address.value.length < 41) {
    alert("名字和钱包不能为空");
    return;
  }
  let groupVal = "-";
  if (group.value) {
    groupVal = group.value;
  }
  api
    .post("/api/wallets", {
      name: name.value,
      address: address.value,
      group: groupVal,
    })
    .then((res) => {
      const data = res.data;
      if (data.error) {
        alert(data.error);
      } else {
        address.value = "";
      }
    });
};

defineExpose({ open });
</script>
