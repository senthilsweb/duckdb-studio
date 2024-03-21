<template>
   <div class="py-5">
    <div class="container max-w-xl   mx-auto">
      <img src="/404-error-01.jpg" alt="404 error" />
    </div>
  </div>
  <!-- typeof error.statusCode: -->
  <!-- string from server-side -->
  <!-- number on client-side -->
  typeof error.statusCode =
  <span v-html="`${typeof error.statusCode}`"> </span>
  <hr />

  <p>
    error.statusCode: {{ error.statusCode }} <br />
    error.message: {{ error.message }}
  </p>
  <hr />

  <!-- customise 404 message from template section -->
  <h3>Custom 404 message</h3>
  <p v-if="error.statusCode === 404 || '404'">
    [template]: Oops page not available 🥺
  </p>

  <!-- Redirect to home page -->
  <button @click="handleError">Clear errors</button>
</template>

<script setup>
// default props available on error.vue
const props = defineProps({
  error: Object,
});

// customise 404 message from script section
const error = useError();


// clear error and redirect to home page
const handleError = () => clearError({ redirect: '/' });
</script>
