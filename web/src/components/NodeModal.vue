<template>
    <div class="modal" id="AddNodeModal" :class="{ 'is-flex' : isOpen }">
        <div class="modal-background"></div>
        <div class="modal-card">
            <header class="modal-card-head">
                <p class="modal-card-title">Add Node</p>
                <button class="delete" aria-label="close" @click="close"></button>
            </header>
            <section class="modal-card-body">
                <div class="field is-horizontal">
                    <div class="field-label is-normal">
                      <label class="label">ID</label>
                    </div>
                    <div class="field-body">
                      <div class="field">
                        <div class="control">
                          <input class="input" type="text" maxlength="8" placeholder="ID" v-model.lazy.trim="id">
                        </div>
                      </div>
                    </div>
                </div>
                <div class="field is-horizontal">
                    <div class="field-label is-normal">
                      <label class="label">IP</label>
                    </div>
                    <div class="field-body">
                      <div class="field">
                        <div class="control">
                          <input class="input" type="text" maxlength="15" placeholder="IP" v-model.lazy.trim="ip">
                        </div>
                      </div>
                    </div>
                </div>
                <div class="field is-horizontal">
                    <div class="field-label is-normal">
                      <label class="label">ISP</label>
                    </div>
                    <div class="field-body">
                      <div class="field">
                        <div class="control">
                          <input class="input" type="text" maxlength="15" placeholder="ISP" v-model.lazy.trim="isp">
                        </div>
                      </div>
                    </div>
                </div>
                <div class="field is-horizontal">
                    <div class="field-label is-normal">
                      <label class="label">Location</label>
                    </div>
                    <div class="field-body">
                      <div class="field">
                        <div class="control">
                          <input class="input" type="text" maxlength="32" placeholder="Location" v-model.lazy.trim="location">
                        </div>
                      </div>
                    </div>
                </div>
                <div class="field is-horizontal">
                    <div class="field-label is-normal">
                      <label class="label">CPU</label>
                    </div>
                    <div class="field-body">
                      <div class="field has-addons">
                        <div class="control">
                          <input class="input" type="number" maxlength="32" placeholder="CPU" v-model.lazy.trim="cpu">
                        </div>
                        <p class="control">
                            <a class="button">Cores</a>
                        </p>
                      </div>
                    </div>
                </div>
                <div class="field is-horizontal">
                    <div class="field-label is-normal">
                      <label class="label">RAM</label>
                    </div>
                    <div class="field-body">
                      <div class="field has-addons">
                        <div class="control">
                          <input class="input" type="number" maxlength="32" placeholder="RAM" v-model.lazy.trim="ram">
                        </div>
                        <p class="control">
                            <a class="button">GB</a>
                        </p>
                      </div>
                    </div>
                </div>
                <div class="field is-horizontal">
                    <div class="field-label is-normal">
                      <label class="label">Disk</label>
                    </div>
                    <div class="field-body">
                      <div class="field has-addons">
                        <div class="control">
                          <input class="input" type="number" maxlength="32" placeholder="Disk" v-model.lazy.trim="disk">
                        </div>
                        <p class="control">
                            <a class="button">GB</a>
                        </p>
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
                          <input class="input" type="number" maxlength="32" placeholder="Bandwidth" v-model.lazy.trim="bandwidth">
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
                          <input class="input" type="number" maxlength="32" placeholder="Traffic" v-model.lazy.trim="traffic">
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
                          <input class="input" type="number" maxlength="32" placeholder="Price" v-model.lazy.trim="price">
                        </div>
                        <p class="control">
                            <a class="button">￥</a>
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
                          <input class="input" type="date" maxlength="32" placeholder="Renew" v-model.lazy.trim="renew">
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
    const isOpen = ref(false)
    const id = ref('')
    const ip = ref('')
    const isp = ref('')
    const location = ref('')
    const cpu = ref()
    const ram = ref()
    const disk = ref()
    const bandwidth = ref('')
    const traffic = ref('')
    const price = ref()
    const renew = ref()

    const open = () => {
        isOpen.value = true
    }

    const close = () => {
        isOpen.value = false
    }

    const add = () => {
        if(!ip.value){
            return;
        }
        console.log(id.value)
        console.log(ip.value)
        console.log(isp.value)
        console.log(location.value)
    }

    defineExpose({open})
</script>