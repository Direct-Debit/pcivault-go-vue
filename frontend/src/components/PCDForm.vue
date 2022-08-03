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

  // If you want to change the way submit button looks, you can select it by id and just update it with raw javascript.
  // Uncomment the code below to see how it works
  /*
   let submitButton = document.getElementById("pcivault-pcd-form-button-submit")
   submitButton.style.background = "blue"
   submitButton.innerHTML = "Submit"

   // you can even add and remove classes to completely override the default styling.
   submitButton.classList.remove("card-form__button")
   submitButton.classList.add("custom_style")
   */
})
</script>

<script>
export default {
  name: "PCDForm",
}
</script>

<style>
.custom_style {
  padding: 25px;
  color: red;
}
</style>
