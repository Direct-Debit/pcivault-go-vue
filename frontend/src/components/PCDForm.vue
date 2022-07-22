<template>
  <div id="pcd_form"></div>
</template>

<script setup>
import {onMounted} from 'vue'

onMounted(async () => {
  let res = await fetch("/backend/get-secret")
  let submit_info = await res.json()

  window.pcd_form(document.getElementById("pcd_form"), {
    submit_secret: submit_info["secret"],
    submit_url: submit_info["url"],
  })
})
</script>

<script>
export default {
  name: "PCDForm",
}
</script>
