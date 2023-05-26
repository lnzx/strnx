<template>
  <div class="modal" id="terminalModal" :class="{ 'is-active': isOpen }">
    <div class="modal-background"></div>
    <div class="modal-card" style="width: 60%">
      <header class="modal-card-head">
        <p class="modal-card-title">Terminal {{ nameRef }} ({{ ipRef }})</p>
        <button class="delete" aria-label="close" @click="close"></button>
      </header>
      <section class="modal-card-body" style="padding: 0 0 15px 15px">
        <div id="term" style="height: 100%"></div>
      </section>
    </div>
  </div>
</template>

<script setup>
import {
  createTerminal,
  initTerminalSession,
  sleep,
} from "@/composable/useTerminal";

const props = defineProps({
  name: String,
  ip: String,
});

const isOpen = ref(false);
const nameRef = ref("");
const ipRef = ref("");

let term;

const open = async (name, ip) => {
  nameRef.value = name;
  ipRef.value = ip;
  isOpen.value = true;

  term = createTerminal("term", ip);
  console.log("create terminal", term);

  await initTerminalSession(term, ip);
};

const close = () => {
  isOpen.value = false;
  if (term) {
    console.log("dispose terminal");
    term.dispose();
  }
};

defineExpose({ open });
</script>
