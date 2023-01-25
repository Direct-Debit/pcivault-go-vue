<template>
  <div v-for="token in tokens.tokens">
    <button @click="() => stripe(token)">
      {{token.token}}
    </button>
  </div>
</template>

<script>
export default {
  name: "TokenList"
}
</script>

<script setup>
import {onMounted} from 'vue'
import {reactive} from 'vue'

const tokens = reactive({tokens: []})

onMounted(async () => {
  let res = await fetch("/backend/get-tokens")
  let res_json = await res.json()
  tokens.tokens = res_json["tokens"]
})

function stripe(token) {
  fetch("/backend/stripe", {body: JSON.stringify(token), method:"POST"})
}
</script>

<style scoped>

</style>